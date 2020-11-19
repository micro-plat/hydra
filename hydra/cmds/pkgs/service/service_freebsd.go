// Copyright 2019 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

package service

import (
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

const version = "freebsd"

type freebsdSystem struct{}

func (freebsdSystem) String() string {
	return version
}
func (freebsdSystem) Detect() bool {
	return true
}
func (freebsdSystem) Interactive() bool {
	return interactive
}
func (freebsdSystem) New(i Interface, c *Config) (Service, error) {
	s := &freebsdService{
		i:      i,
		Config: c,
	}

	return s, nil
}

func init() {
	ChooseSystem(freebsdSystem{})
}

var interactive = false

func init() {
	var err error
	interactive, err = isInteractive()
	if err != nil {
		panic(err)
	}
}

func isInteractive() (bool, error) {
	return os.Getenv("IS_DAEMON") != "1", nil
}

type freebsdService struct {
	i Interface
	*Config
}

func (s *freebsdService) String() string {
	if len(s.DisplayName) > 0 {
		return s.DisplayName
	}
	return s.Name
}

func (s *freebsdService) Platform() string {
	return version
}

func (s *freebsdService) template() *template.Template {
	functions := template.FuncMap{
		"bool": func(v bool) string {
			if v {
				return "true"
			}
			return "false"
		},
	}

	customConfig := s.Option.string(optionSysvScript, "")

	if customConfig != "" {
		return template.Must(template.New("").Funcs(functions).Parse(customConfig))
	} else {
		return template.Must(template.New("").Funcs(functions).Parse(rcScript))
	}
}

func (s *freebsdService) configPath() (cp string, err error) {
	cp = "/usr/local/etc/rc.d/" + s.Config.Name
	return
}

func (s *freebsdService) Install() error {
	path, err := s.execPath()
	if err != nil {
		return err
	}

	// write start script
	confPath, err := s.configPath()
	if err != nil {
		return err
	}
	_, err = os.Stat(confPath)
	if err == nil {
		return ErrHasInstalled
	}

	f, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var to = &struct {
		*Config
		Path string
	}{
		s.Config,
		path,
	}

	err = s.template().Execute(f, to)
	if err != nil {
		return err
	}

	if err = os.Chmod(confPath, 0755); err != nil {
		return err
	}

	return nil
}

func (s *freebsdService) Uninstall() error {
	confPath, err := s.configPath()
	if err != nil {
		return err
	}
	_, err = os.Stat(confPath)
	if err != nil {
		return ErrNotInstalled
	}
	if err := os.Remove(confPath); err != nil {
		return err
	}
	return nil
}

func (s *freebsdService) Status() (Status, error) {
	cp, err := s.configPath()
	if err != nil {
		return StatusUnknown, err
	}

	if _, err = os.Stat(cp); os.IsNotExist(err) {
		return StatusStopped, ErrNotInstalled
	}

	status, _, err := runCommand("service", false, s.Name, "status")
	if status == 1 {
		return StatusStopped, nil
	} else if err != nil {
		return StatusUnknown, err
	}
	return StatusRunning, nil
}

func (s *freebsdService) Start() error {
	if !s.isInstalled() {
		return ErrNotInstalled
	}
	status, err := s.Status()
	if err != nil {
		return err
	}
	if status == StatusRunning {
		return ErrIsRunning
	}
	return run("service", s.Name, "start")
}

func (s *freebsdService) Stop() error {
	if !s.isInstalled() {
		return ErrNotInstalled
	}
	status, err := s.Status()
	if err != nil {
		return err
	}
	if status == StatusStopped {
		return ErrHasStopped
	}
	return run("service", s.Name, "stop")
}

func (s *freebsdService) Restart() error {
	return run("service", s.Name, "restart")
}

func (s *freebsdService) Run() error {
	var err error

	err = s.i.Start(s)
	if err != nil {
		return err
	}

	s.Option.funcSingle(optionRunWait, func() {
		var sigChan = make(chan os.Signal, 3)
		signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)
		<-sigChan
	})()

	return s.i.Stop(s)
}

func (s *freebsdService) Logger(errs chan<- error) (Logger, error) {
	if interactive {
		return ConsoleLogger, nil
	}
	return s.SystemLogger(errs)
}

func (s *freebsdService) SystemLogger(errs chan<- error) (Logger, error) {
	return newSysLogger(s.Name, errs)
}


func (s *freebsdService) isInstalled() bool {
	confPath, err := s.configPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(confPath)
	if err != nil {
		return false
	}
	return true
}

var rcScript = `#!/bin/sh

# PROVIDE: {{.Name}}
# REQUIRE: SERVERS
# KEYWORD: shutdown

. /etc/rc.subr

name="{{.Name}}"
{{.Name}}_env="IS_DAEMON=1"
pidfile="/var/run/${name}.pid"
command="/usr/sbin/daemon"
daemon_args="-P ${pidfile} -r -t \"${name}: daemon\"{{if .WorkingDirectory}} -c {{.WorkingDirectory}}{{end}}"
command_args="${daemon_args} {{.Path}}{{range .Arguments}} {{.}}{{end}}"

run_rc_command "$1"
`
