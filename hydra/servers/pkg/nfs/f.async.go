package nfs

import "time"

//async 异步处理任务
type async struct {
	remoting       *remoting
	local          *local
	reportChan     chan eFileFPLists
	downloadChan   chan *eFileFP
	queryChan      chan struct{}
	reportExit     bool
	downloadExit   bool
	maxGocoroutine int
}

func newAsync(l *local, r *remoting) *async {
	m := &async{
		remoting:       r,
		local:          l,
		reportChan:     make(chan eFileFPLists, 10000),
		downloadChan:   make(chan *eFileFP, 10000),
		queryChan:      make(chan struct{}, 1),
		maxGocoroutine: 10,
	}
	go m.loopReport()
	go m.loopQuery()
	for i := 0; i < m.maxGocoroutine; i++ {
		go m.loopDownload()
	}
	return m
}

//DoReport 上报指纹信息
func (m *async) DoReport(f eFileFPLists) {
	if len(f) == 0 {
		return
	}
	m.reportChan <- f
}

//DoQuery 执行远程服务查询
func (m *async) DoQuery() {
	select {
	case m.queryChan <- struct{}{}:
	default:
	}
}

//DoDownload 下载文件
func (m *async) DoDownload(f *eFileFP) {
	m.downloadChan <- f
}

//loopReport 循环处理上报
func (m *async) loopReport() {
	for {
		select {
		case f, ok := <-m.reportChan:
			if !ok {
				return
			}
			m.remoting.Report(f)
		}
	}
}

//loopReport 循环处理上报
func (m *async) loopQuery() {
	for {
		select {
		case _, ok := <-m.queryChan:
			if !ok {
				return
			}

			//查询远程服务器列表
			mp, err := m.remoting.Query()
			if err != nil {
				go func() {
					time.Sleep(time.Second * 60)
					m.DoQuery()
				}()
				continue
			}
			//结合外部传入，与当前服务器，进行整体合并，并进行通知
			mp.Merge(m.local.GetFPs())
			m.DoReport(mp) //?应该包含自己
		}
	}
}

//loopDownload  循环处理下载
func (m *async) loopDownload() {
	for {
		select {
		case f, ok := <-m.downloadChan:
			if !ok {
				return
			}
			if m.local.Has(f.Path) {
				continue
			}
			//从远程拉取文件
			buff, err := m.remoting.Pull(f)
			if err != nil {
				go func() {
					time.Sleep(time.Second * 60)
					m.DoDownload(f) //出错自动下载
				}()
				continue
			}
			fx, err := m.local.SaveFile(f.Path, buff, f.Hosts...)
			if err != nil {
				go func() {
					time.Sleep(time.Second * 60)
					m.DoDownload(f) //出错自动下载
				}()
				continue
			}
			m.DoReport(fx.GetMAP())
		}
	}
}

//Close 关闭任务
func (m *async) Close() error {
	close(m.reportChan)
	close(m.downloadChan)
	close(m.queryChan)
	return nil
}
