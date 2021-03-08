package zk

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"sync"

	"github.com/micro-plat/lib4go/logger"
	"github.com/samuel/go-zookeeper/zk"
)

// TIMEOUT 连接zk服务器操作的超时时间
var TIMEOUT = time.Second

/*
type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Printf(string, ...interface{})
}
*/
var (
	ErrColientCouldNotConnect = errors.New("zk: could not connect to the server")
	ErrClientConnClosing      = errors.New("zk: the client connection is closing")
)

//ZookeeperClient zookeeper客户端
type ZookeeperClient struct {
	servers   []string
	timeout   time.Duration
	clock     sync.Mutex
	conn      *zk.Conn
	eventChan <-chan zk.Event
	Log       logger.ILogging
	useCount  int32
	isConnect bool
	once      sync.Once
	CloseCh   chan struct{}
	ACL       []zk.ACL
	digest    bool
	userName  string
	password  string
	// 是否是手动关闭
	done bool
}

//New 连接到Zookeeper服务器
func New(servers []string, timeout time.Duration, opts ...Option) (*ZookeeperClient, error) {
	log := logger.GetSession("zk", logger.CreateSession())
	return NewWithLogger(servers, timeout, log, opts...)
}

//NewWithLogger 连接到Zookeeper服务器
func NewWithLogger(servers []string, timeout time.Duration, logger logger.ILogging, opts ...Option) (*ZookeeperClient, error) {
	client := &ZookeeperClient{servers: servers, timeout: timeout, useCount: 0}
	client.CloseCh = make(chan struct{})
	client.Log = logger
	for _, opt := range opts {
		opt(client)
	}
	if client.digest {
		client.ACL = zk.DigestACL(zk.PermAll, client.userName, client.password)
	} else {
		client.ACL = zk.WorldACL(zk.PermAll)
	}
	return client, nil
}

//Connect 连接到远程zookeeper服务器
func (client *ZookeeperClient) Connect() (err error) {
	if client.conn == nil {
		conn, eventChan, err := zk.Connect(client.servers, client.timeout)
		if err != nil {
			return err
		}
		if client.digest {
			if err := conn.AddAuth("digest",
				[]byte(fmt.Sprintf("%s:%s",
					client.userName,
					client.password))); err != nil {
				return nil
			}
		}
		client.conn = conn
		client.conn.SetLogger(client.Log)
		client.eventChan = eventChan
		go client.eventWatch()
	}
	atomic.AddInt32(&client.useCount, 1)
	for client.conn.State() != zk.StateHasSession {
		time.Sleep(50 * time.Millisecond)
	}
	// time.Sleep(time.Second)
	client.isConnect = true
	return
}

//IsConnected 是否已连接到服务器
func (client *ZookeeperClient) IsConnected() bool {
	return client.isConnect
}

//Reconnect 重新连接服务器
func (client *ZookeeperClient) Reconnect() (err error) {
	client.isConnect = false
	if client.conn != nil {
		client.conn.Close()
	}
	client.done = false
	return client.Connect()
}

//CanWirteDataInDir 目录中能否写入数据
func (client *ZookeeperClient) CanWirteDataInDir() bool {
	return true
}

//Close 关闭服务器
func (client *ZookeeperClient) Close() error {
	atomic.AddInt32(&client.useCount, -1)
	if client.useCount > 0 {
		return nil
	}

	if client.conn != nil {
		client.once.Do(client.conn.Close)
	}

	client.isConnect = false
	client.done = true
	client.once.Do(func() {
		close(client.CloseCh)
	})
	return nil
}

func (client *ZookeeperClient) GetSeparator() string {
	return "/"
}

var baseVal int64

func init() {
	baseTime := time.Date(time.Now().Year()-10, 1, 1, 0, 0, 0, 0, time.Local)
	baseVal = time.Now().Sub(baseTime).Nanoseconds() / 1e6
}

func getVersion(stat *zk.Stat) int32 {
	if stat == nil {
		return 0
	}
	curtime := stat.Mtime
	return int32((curtime - baseVal) / 1e3)
}
