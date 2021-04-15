package nfs

import (
	"io/fs"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (l *module) watch() {
	if !l.c.Watch {
		return
	}
	var err error
	l.fsWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		trace("不支持文件监控", err)
		return
	}
	filepath.WalkDir(l.local.path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || l.local.exclude(d.Name()) {
			return nil
		}
		l.fsWatcher.Add(path)
		return nil
	})
	go l.delayCheck()
	for {
		select {
		case ev, ok := <-l.fsWatcher.Events:
			if !ok || l.done {
				return
			}
			//监控文件夹变动
			if ev.Op&fsnotify.Create == fsnotify.Create ||
				ev.Op&fsnotify.Write == fsnotify.Write ||
				ev.Op&fsnotify.Remove == fsnotify.Remove ||
				ev.Op&fsnotify.Rename == fsnotify.Rename {
				if l.done {
					return
				}
				_, name := filepath.Split(ev.Name)
				if l.local.exclude(name) {
					continue
				}
				select {
				case l.checkChan <- struct{}{}:
				default:
				}
			}
		}
	}
}
func (l *module) delayCheck() {
	tk := time.Tick(time.Second * 10)
	for {
		select {
		case <-tk:
			select {
			case _, ok := <-l.checkChan:
				if !ok || l.done {
					return
				}
				if l.local.FindChange() {
					l.async.DoReport(l.local.GetFPs())
				}
			default:
			}
		}
	}
}
