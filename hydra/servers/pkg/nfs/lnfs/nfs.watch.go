package lnfs

import (
	"io/fs"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (l *Module) watch() {
	if !l.c.Watch {
		return
	}
	var err error
	l.fsWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		trace().error("不支持文件监控", err)
		return
	}
	filepath.WalkDir(l.Local.path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || l.Local.exclude(path, d.Name()) {
			return nil
		}
		l.fsWatcher.Add(path)
		return nil
	})
	go l.delayCheck()
	for {
		select {
		case ev, ok := <-l.fsWatcher.Events:
			if !ok || l.Done {
				return
			}
			//监控文件夹变动
			if ev.Op&fsnotify.Create == fsnotify.Create ||
				ev.Op&fsnotify.Write == fsnotify.Write ||
				ev.Op&fsnotify.Remove == fsnotify.Remove ||
				ev.Op&fsnotify.Rename == fsnotify.Rename {
				if l.Done {
					return
				}
				_, name := filepath.Split(ev.Name)
				if l.Local.exclude(ev.Name, name) {
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
func (l *Module) delayCheck() {
	tk := time.Tick(time.Second * 10)
	for {
		select {
		case <-tk:
			select {
			case _, ok := <-l.checkChan:
				if !ok || l.Done {
					return
				}
				if l.Local.FindChange() {
					l.async.DoReport(l.Local.GetFPs())
				}
			default:
			}
		}
	}
}
