package pipes

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RunString Convert a shell command with a series of pipes into
// correspondingly piped list of *exec.Cmd
// If an arg has spaces, this will fail
func RunString(s string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	sp := strings.Split(s, "|")
	cmds := make([]*exec.Cmd, len(sp))
	// create the commands
	for i, c := range sp {
		cs := strings.Split(strings.TrimSpace(c), " ")
		cmd := cmdFromStrings(cs)
		cmds[i] = cmd
	}

	cmds = AssemblePipes(cmds, nil, buf)
	if err := RunCmds(cmds); err != nil {
		return "", err
	}

	b := buf.Bytes()
	return string(b), nil
}

func cmdFromStrings(cs []string) *exec.Cmd {
	if len(cs) == 1 {
		return exec.Command(cs[0])
	} else if len(cs) == 2 {
		return exec.Command(cs[0], cs[1])
	}
	return exec.Command(cs[0], cs[1:]...)
}

//RunStrings  Convert sequence of tokens into commands,
// using "|" as a delimiter
func RunStrings(tokens ...string) (string, error) {
	if len(tokens) == 0 {
		return "", nil
	}
	buf := bytes.NewBuffer([]byte{})
	cmds := []*exec.Cmd{}
	args := []string{}
	// accumulate tokens until a |
	for _, t := range tokens {
		if t != "|" {
			args = append(args, t)
		} else {
			cmds = append(cmds, cmdFromStrings(args))
			args = []string{}
		}
	}
	cmds = append(cmds, cmdFromStrings(args))
	cmds = AssemblePipes(cmds, nil, buf)
	if err := RunCmds(cmds); err != nil {
		return "", fmt.Errorf("%s; %s", err.Error(), string(buf.Bytes()))
	}

	b := buf.Bytes()
	return string(b), nil
}

//AssemblePipes  Pipe stdout of each command into stdin of next
func AssemblePipes(cmds []*exec.Cmd, stdin io.Reader, stdout io.Writer) []*exec.Cmd {
	cmds[0].Stdin = stdin
	cmds[0].Stderr = stdout
	// assemble pipes
	for i, c := range cmds {
		if i < len(cmds)-1 {
			cmds[i+1].Stdin, _ = c.StdoutPipe()
			cmds[i+1].Stderr = stdout
		} else {
			c.Stdout = stdout
			c.Stderr = stdout
		}
	}
	return cmds
}

// RunCmds run series of piped commands
func RunCmds(cmds []*exec.Cmd) error {
	// start processes in descending order
	for i := len(cmds) - 1; i > 0; i-- {
		if err := cmds[i].Start(); err != nil {
			return err
		}
	}
	// run the first process
	if err := cmds[0].Run(); err != nil {
		return err
	}
	// wait on processes in ascending order
	for i := 1; i < len(cmds); i++ {
		if err := cmds[i].Wait(); err != nil {
			return err
		}
	}
	return nil
}

//BashRun 通过bash 执行命令
func BashRun(cmd string) (content string, err error) {
	cmd1 := exec.Command("/bin/bash", "-c", cmd)
	var out bytes.Buffer
	cmd1.Stdout = &out
	cmd1.Stderr = &out
	err = cmd1.Start()
	if err != nil {
		return
	}
	err = cmd1.Wait()
	if err != nil {
		return
	}

	content = strings.Trim(out.String(), "\n")
	return
}

//Run 通过管道方式执行命令
func Run(cmds []*exec.Cmd) (content string, err error) {
	var out bytes.Buffer
	AssemblePipes(cmds, os.Stdin, &out)
	err = RunCmds(cmds)
	time.Sleep(time.Second)
	if err != nil {
		err = fmt.Errorf("命令执行失败:%s,err:%v", strings.Trim(out.String(), "\n"), err)
		return
	}
	content = strings.Trim(out.String(), "\n")
	return
}
