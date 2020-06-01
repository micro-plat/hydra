package influxdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/lib4go/transform"
)

// ConfigOptions influxdb配置
type ConfigOptions struct {
	Address   string `json:"address"`
	DbName    string `json:"db"`
	UserName  string `json:"user"`
	Password  string `json:"password"`
	RowFormat string `json:"row"`
}

// InfluxDB 上下文
type InfluxDB struct {
	config ConfigOptions
}

// New 新建一个influxdb的环境
func New(config ConfigOptions) (i *InfluxDB, err error) {
	i = &InfluxDB{}
	i.config = config
	if strings.EqualFold(i.config.Address, "") ||
		strings.EqualFold(i.config.DbName, "") ||
		strings.EqualFold(i.config.RowFormat, "") {
		err = errors.New("ConfigOptions必须参数不能为空")
		return
	}
	return
}

// SaveString 保存json格式的字符串
func (db *InfluxDB) SaveString(rows string) (err error) {
	var data []map[string]interface{}
	err = json.Unmarshal([]byte(rows), &data)
	if err != nil {
		return fmt.Errorf("influxdb SaveString 反序列化字符串失败:%v", err)
	}
	return db.Save(data)
}

// Save 保存map类型的数据
func (db *InfluxDB) Save(rows []map[string]interface{}) (err error) {
	url := fmt.Sprintf("%s/write?db=%s", db.config.Address, db.config.DbName)
	var datas []string
	for i := 0; i < len(rows); i++ {
		d := transform.NewMaps(rows[i])
		datas = append(datas, d.Translate(db.config.RowFormat))
	}
	data := strings.Join(datas, "\n")
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("influxdb Post fail:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}
	err = fmt.Errorf("influxdb save error:%d", resp.StatusCode)
	return
}
