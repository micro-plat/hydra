package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/archiver"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

// FileAppender 文件FileAppender
type FileAppender struct {
	writers    cmap.ConcurrentMap
	ticker     *time.Ticker
	closNotify chan struct{}
	done       chan struct{}
}

// NewFileAppender 构建file FileAppender
func NewFileAppender() *FileAppender {
	a := &FileAppender{
		done:    make(chan struct{}),
		writers: cmap.New(6),
		ticker:  time.NewTicker(time.Minute * 1),
	}
	go a.clean()
	return a
}
func (a *FileAppender) Write(layout *Layout, event *LogEvent) error {
	p := event.Transform(layout.Path)
	_, w, err := a.writers.SetIfAbsentCb(p, func(input ...interface{}) (interface{}, error) {
		return newWriter(p, layout)
	})
	if err != nil {
		return err
	}
	w.(*writer).Write(event)
	return nil
}
func (a *FileAppender) clean() {
	for {
		select {
		case <-a.done:
			return
		case <-a.ticker.C:
			a.writers.RemoveIterCb(func(key string, value interface{}) bool {
				w := value.(*writer)
				if time.Since(w.lastWrite) > 5*time.Minute {
					w.Write(GetEndWriteEvent()) //向日志发送结速写入事件
					w.Close()                   //等待所有日志被写入文件
					//如果文件是昨天的，启动跨天处理
					current := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
					if w.lastWrite.Before(current) {
						go afterDay(w.layout, key)
					}
					return true
				}
				return false
			})
		}
	}

}
func afterDay(layout *Layout, path string) {
	if layout.AutoCompress {
		err := archiver.Zip.Make(fmt.Sprintf("%s.zip", path), []string{path})
		if err != nil {
			SysLog.Errorf("压缩文件：%s出现异常:%v", path, err)
			return
		}
		err = os.Remove(path)
		if err != nil {
			SysLog.Errorf("移除文件:%s出现异常:%v", path, err)
		}
	}
	if layout.MaxDays > 0 {
		current := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		dt := current.AddDate(0, 0, -1*layout.MaxDays)
		files, err := getFiles(path, dt)
		if err != nil {
			return
		}
		for _, f := range files {
			if !strings.HasPrefix(strings.ToLower(f), ".zip") {
				continue
			}
			err := os.Remove(f)
			if err != nil {
				SysLog.Errorf("移除文件:%s出现异常:%v", path, err)
			}
		}
	}
}
func getFiles(path string, dt time.Time) (files []string, err error) {
	// 确定目标目录路径（处理传入路径可能是文件的情况）
	targetDir := path
	if fileInfo, err := os.Stat(path); err == nil {
		if !fileInfo.IsDir() {
			targetDir = filepath.Dir(path) // 如果是文件则取其所在目录[[7,12]]
		}
	} else {
		return nil, err
	}

	// 读取目录下的所有条目
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, err
	}

	// 遍历目录条目
	for _, entry := range entries {
		if entry.IsDir() {
			continue // 跳过子目录[[1,4]]
		}

		// 获取文件信息
		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的文件
		}

		// 检查修改时间是否早于指定时间
		if info.ModTime().Before(dt) {
			fullPath := filepath.Join(targetDir, entry.Name())
			files = append(files, fullPath)
		}
	}

	return files, nil
}

// Close 关闭组件
func (a *FileAppender) Close() error {
	close(a.done)
	a.writers.RemoveIterCb(func(key string, w interface{}) bool {
		w.(*writer).Write(GetEndWriteEvent())
		w.(*writer).Close()
		return true
	})
	return nil
}
