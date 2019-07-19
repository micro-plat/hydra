package cron

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func (s *Processor) saveHistory(t iCronTask) error {
	if s.redisClient == nil {
		return nil
	}
	data, err := t.GetTaskExecutionRecord()
	if err != nil {
		return err
	}
	member := redis.Z{Member: data}
	member.Score, _ = strconv.ParseFloat(time.Now().Format("20060102150405"), 64)

	key := fmt.Sprintf(s.historyNode, t.GetName())
	r, err := s.redisClient.ZAdd(key, member).Result()
	if err != nil || r == 0 {
		err = fmt.Errorf("保存cron执行记录失败:c:%d,err:%v", r, err)
		return err
	}
	min := time.Now().Add(time.Hour * -3600).Format("20060102150405")
	max := time.Now().Add(time.Hour * -360).Format("20060102150405") //删除15天前的数据
	_, err = s.redisClient.ZRemRangeByScore(key, min, max).Result()
	return err
}
