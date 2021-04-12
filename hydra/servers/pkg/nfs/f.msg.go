package nfs

import "sync"

type msg struct {
	r            *remoting
	l            *local
	hosts        []string
	reportChan   chan map[string]eFileFPLists
	downloadChan chan *eFileFP
	reportExit   bool
	downloadExit bool
	once         sync.Once
}

func newMsg(l *local, r *remoting) *msg {
	m := &msg{
		r:            r,
		l:            l,
		reportChan:   make(chan map[string]eFileFPLists, 1000),
		downloadChan: make(chan *eFileFP, 1000),
	}
	go m.loopReport()
	go m.loopDownload()
	return m
}
func (m *msg) Update(hosts []string) {
	m.hosts = hosts
}

//Report 上报指纹信息
func (m *msg) Report(f map[string]eFileFPLists) {
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
			m.r.Report(f)
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
			buff, err := m.r.Pull(f)
			if err != nil {
				continue
			}
			fx, err := m.l.SaveFile(f.Path, buff, f.Hosts...)
			if err != nil {
				continue
			}
			m.Report(GetNotify(fx.GetMAP(), m.hosts))
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
