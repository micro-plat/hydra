package nfs

import "sync"

type msg struct {
	reportChan   chan *eFileFP
	downloadChan chan *eFileFP
	r            *remoting
	l            *local
	reportExit   bool
	downloadExit bool
	once         sync.Once
}

func newMsg() *msg {
	m := &msg{
		reportChan:   make(chan *eFileFP, 1000),
		downloadChan: make(chan *eFileFP, 1000),
	}
	go m.loopReport()
	go m.loopDownload()
	return m
}

//Report 上报指纹信息
func (m *msg) Report(f *eFileFP) {
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
			if err != nil {
				return
			}
			fx, err := m.l.SaveFile(f.Path, buff)
			if err != nil {
				return
			}
			m.Report(fx)
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
