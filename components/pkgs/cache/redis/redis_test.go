package redis

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkRedisGet(b *testing.B) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		b.Errorf("获取redis对象失败,err:%+v", err)
		return
	}
	if err = redisClient.Add("taosytest123", "123456", 0); err != nil {
		b.Errorf("添加redis数据失败,err:%+v", err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss, err := redisClient.Get("taosytest123")
		if err != nil {
			b.Errorf("redis中读取数据异常,err:%+v", err)
			return
		}

		if ss != "123456" {
			b.Errorf("redis中数据不匹配异常,ss:%s", ss)
			return
		}
	}
	b.StopTimer()

	if err = redisClient.Delete("taosytest123"); err != nil {
		b.Errorf("删除redis数据失败,err:%+v", err)
		return
	}
}

func TestRedisAdd(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}
	if err = redisClient.Add("taosytest123", "123456", 0); err != nil {
		t.Errorf("添加redis数据失败,err:%+v", err)
		return
	}
	ss, err := redisClient.Get("taosytest123")
	if err != nil {
		t.Errorf("redis中读取数据异常,err:%+v", err)
		return
	}

	if ss != "123456" {
		t.Errorf("redis中数据不匹配异常,ss:%s", ss)
		return
	}

	if err = redisClient.Add("taosytest123", "123456", 0); err == nil {
		t.Error("添加redis数据不能成功")
		return
	}

	if err = redisClient.Delete("taosytest123"); err != nil {
		t.Errorf("删除redis数据失败,err:%+v", err)
		return
	}
}

func TestRedisSet(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	if err = redisClient.Set("taosytest123", "123456", 0); err != nil {
		t.Errorf("添加redis数据失败,err:%+v", err)
		return
	}

	ss, err := redisClient.Get("taosytest123")
	if err != nil {
		t.Errorf("redis中读取数据异常,err:%+v", err)
		return
	}

	if ss != "123456" {
		t.Errorf("redis中数据不匹配异常,ss:%s", ss)
		return
	}

	if err = redisClient.Set("taosytest123", "654321", 0); err != nil {
		t.Errorf("添加redis数据失败11,err:%+v", err)
		return
	}

	ss, err = redisClient.Get("taosytest123")
	if err != nil {
		t.Errorf("redis中读取数据异常,err:%+v", err)
		return
	}

	if ss != "654321" {
		t.Errorf("redis中数据不匹配异常,ss:%s", ss)
		return
	}

	if err = redisClient.Delete("taosytest123"); err != nil {
		t.Errorf("删除redis数据失败,err:%+v", err)
		return
	}
}

func TestRedisX(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	n, err := redisClient.Increment("taosytest123", 1)
	if err != nil {
		t.Errorf("添加redis数据失败,err:%+v", err)
		return
	}

	t.Errorf("n:%d \n", n)
	res, err := redisClient.Get("taosytest123")
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	t.Errorf("res:%s \n", res)

}

func TestRedisY(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	redisClient.Delete("taosytest123")
	var add int64 = 0
	var deduct int64 = 0
	for i := 0; i < 100; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				n := time.Now().UnixNano() % 2
				if n == 0 {
					if _, err := redisClient.Decrement("taosytest123", 1); err != nil {
						t.Errorf("添加redis数据失败,err:%+v", err)
					} else {
						deduct++
					}
				} else {
					if _, err := redisClient.Increment("taosytest123", 1); err != nil {
						t.Errorf("添加redis数据失败,err:%+v", err)
					} else {
						add++
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	t.Errorf("add:%d,deduct:%d", add, deduct)
	res, err := redisClient.Get("taosytest123")
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	r, _ := strconv.ParseInt(res, 10, 64)
	if add-deduct != r {
		t.Errorf("11111111111111,res:%s,add-deduct:%d", res, add-deduct)
		return
	}

	if err = redisClient.Delete("taosytest123"); err != nil {
		t.Errorf("删除redis数据失败,err:%+v", err)
		return
	}
}

func TestXXX(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	list := []string{}
	c := uint64(0)
	for {
		res := redisClient.client.Scan(c, "17sup_v2_debug:inbound_list_main:*", 100)
		arry, c1, err := res.Result()
		if err != nil {
			t.Errorf("添加redis数据失败,err:%+v", err)
			return
		}
		if arry != nil && len(arry) > 0 {
			allArry := make([]string, len(arry)+len(list))
			copy(allArry, list)
			copy(allArry[len(list):], arry)
			list = allArry
		}
		if c1 == 0 {
			break
		}
		c = c1
	}

	t.Errorf("n:%s\n", list)
}

func TestXXX1(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	res := redisClient.client.Keys("17sup_v2_debug:inbound_list_main:*")
	arry, err := res.Result()
	if err != nil {
		t.Errorf("添加redis数据失败,err:%+v", err)
		return
	}

	t.Errorf("n:%s\n", arry)
}

func TestXXX2(t *testing.T) {
	redisClient, err := New([]string{"192.168.0.116:6379"}, `{"addrs":["192.168.0.116:6379","192.168.0.113:6379"],"db":1}`)
	if err != nil {
		t.Errorf("获取redis对象失败,err:%+v", err)
		return
	}

	res, err := redisClient.client.Del("sdfdfsf").Result()

	// res := redisClient.client.Exists("17sup_v2_debug:inbound_list_main:sdddf")
	// arry, err := res.Result()

	t.Errorf("n:%d err:%+v\n", res, err)
}
