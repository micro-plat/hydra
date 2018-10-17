package influxdb

import (
	"errors"
	"fmt"
	"log"
	uurl "net/url"
	"time"
	//"github.com/micro-plat/lib4go/influxdb"
	//"github.com/micro-plat/lib4go/influxdb"
)

type IInfluxClient interface {
	QueryResponse(sql string) (response *Response, err error)
	QueryMaps(sql string) (rx [][]map[string]interface{}, err error)
	Query(sql string) (result string, err error)
	SendLineProto(data string) error
	Send(measurement string, tags map[string]string, fileds map[string]interface{}) error
	Close() error
}

type InfluxClient struct {
	interval time.Duration
	url      uurl.URL
	database string
	username string
	password string
	client   *Client
	closeCh  chan struct{}
	done     bool
}

// newInfluxClient starts a InfluxDB reporter which will post the metrics from the given registry at each d interval with the specified tags
func NewInfluxClient(url, database, username, password string) (*InfluxClient, error) {
	u, err := uurl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse InfluxDB url %s. err=%v", url, err)
	}

	rep := &InfluxClient{
		url:      *u,
		database: database,
		username: username,
		password: password,
		closeCh:  make(chan struct{}),
	}

	if err := rep.makeClient(); err != nil {
		return nil, fmt.Errorf("unable to make InfluxDB client. err=%v", err)
	}
	go rep.run()
	return rep, nil
}

func (r *InfluxClient) makeClient() (err error) {
	r.client, err = NewClient(Config{
		URL:       r.url,
		Timeout:   time.Second * 3,
		UserAgent: "hydra",
		Username:  r.username,
		Password:  r.password,
	})
	return
}

func (r *InfluxClient) run() {
	pingTicker := time.Tick(time.Second * 5)
	for {
		select {
		case <-r.closeCh:
			r.client = nil
			return
		case <-pingTicker:
			_, _, err := r.client.Ping()
			if err != nil {
				log.Printf("got error while sending a ping to InfluxDB, trying to recreate client. err=%v", err)
				if err = r.makeClient(); err != nil {
					log.Printf("unable to make InfluxDB client. err=%v", err)
				}
			}
		}
	}
}
func (r *InfluxClient) QueryResponse(sql string) (response *Response, err error) {
	if r.done {
		return nil, errors.New("连接已关闭")
	}
	response, err = r.client.Query(Query{Command: sql, Database: r.database})
	if err != nil {
		err = fmt.Errorf("query.error:%v", err)
		return
	}
	return
}
func (r *InfluxClient) QueryMaps(sql string) (rx [][]map[string]interface{}, err error) {
	if r.done {
		return nil, errors.New("连接已关闭")
	}
	response, err := r.client.Query(Query{Command: sql, Database: r.database})
	if err != nil {
		err = fmt.Errorf("query.error:%v", err)
		return
	}
	if err = response.Error(); err != nil {
		return nil, fmt.Errorf("response.error:%v", err)
	}
	rx = make([][]map[string]interface{}, 0, len(response.Results))
	for _, v := range response.Results {
		result := make([]map[string]interface{}, 0, 0)
		for _, row := range v.Series {
			for _, value := range row.Values {
				srow := make(map[string]interface{})
				for y, col := range row.Columns {
					srow[col] = value[y]
				}
				result = append(result, srow)
			}

		}
		rx = append(rx, result)
	}
	return rx, nil
}
func (r *InfluxClient) Query(sql string) (result string, err error) {
	if r.done {
		return "", errors.New("连接已关闭")
	}
	response, err := r.client.Query(Query{Command: sql, Database: r.database})
	if err != nil {
		err = fmt.Errorf("query.error:%v", err)
		return
	}
	if err = response.Error(); err != nil {
		return "", fmt.Errorf("response.error:%v", err)
	}
	buf, err := response.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("query.result.marshal.error:%v", err)
		return
	}
	result = string(buf)
	return
}
func (r *InfluxClient) SendLineProto(data string) error {
	if r.done {
		return errors.New("连接已关闭")
	}
	_, err := r.client.WriteLineProtocol(data, r.database, "default", "us", "")
	return err

}
func (r *InfluxClient) Send(measurement string, tags map[string]string, fileds map[string]interface{}) error {
	if r.done {
		return errors.New("连接已关闭")
	}
	var pts []Point
	pts = append(pts, Point{
		Measurement: measurement,
		Tags:        tags,
		Fields:      fileds,
		Time:        time.Now(),
	})

	bps := BatchPoints{
		Points:   pts,
		Database: r.database,
	}

	_, err := r.client.Write(bps)
	return err
}
func (r *InfluxClient) Close() error {
	r.done = true
	close(r.closeCh)
	return nil
}
