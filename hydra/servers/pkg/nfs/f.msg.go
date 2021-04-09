package nfs

import "sync"

type msg struct {
	reportChan   chan eFileFPLists
	downloadChan chan *eFileFP
	r            *remoting
	l            *local
	reportExit   bool
	downloadExit bool
	once         sync.Once
}

func newMsg(l *local, r *remoting) *msg {
	m := &msg{
		r:            r,
		l:            l,
		reportChan:   make(chan eFileFPLists, 1000),
		downloadChan: make(chan *eFileFP, 1000),
	}
	go m.loopReport()
	go m.loopDownload()
	return m
}

//Report 上报指纹信息
func (m *msg) Report(f eFileFPLists) {
	if len(f) == 0 {
		return
	}
	m.reportChan <- f
}

//Download 下载文件
func (m *msg) Download(f *eFileFP) {
	m.downloadChan <- f
}
func (m *msg) loopReport() {
	for {
		select {
		case f, ok := <-m.reportChan:
			if !ok {
				return
			}
			trace("push-1:", f)
			m.r.Push(f)
		}
	}
}
func (m *msg) loopDownload() {
	for {
		select {
		case f, ok := <-m.downloadChan:
			if !ok {
				return
			}
			//从远程拉取文件
			buff, err := m.r.Pull(f.Path, f.Hosts)
			trace("pull-1:", f.Path, err)
			if err != nil {
				trace("pull-1-err:", f.Path, err)
				continue
			}
			fx, err := m.l.SaveFile(f.Path, buff, f.Hosts...)
			trace("pull-2:", f.Path, err)
			if err != nil {
				trace("pull-2-save.err:", f.Path, err)
				continue
			}
			m.Report(fx.GetMAP())
		}
	}
}

func (m *msg) Close() error {
	m.once.Do(func() {
		close(m.reportChan)
		close(m.downloadChan)
	})
	return nil
}
