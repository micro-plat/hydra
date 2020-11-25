package global

import (
	"sort"
	"strings"
)

//IsMaster IsMaster
func IsMaster(master bool, scount int, path string, cldrs []string) (idx int, isMaster bool) {
	if len(cldrs) == 0 {
		return 0, master
	}
	ncldrs := make([]string, 0, len(cldrs))
	for _, v := range cldrs {
		args := strings.SplitN(v, "_", 2)
		ncldrs = append(ncldrs, args[len(args)-1])
	}
	sort.Strings(ncldrs)
	rcount := scount
	if scount == 0 {
		rcount = len(ncldrs)
	}
	index := -1
	for i, v := range ncldrs {
		if strings.HasSuffix(path, v) {
			index = i
			break
		}
	}
	if index == -1 {
		return 0, false
	}
	if scount == 1 {
		return 0, index == 0
	}
	shardingIndex := getSharding(index, rcount)
	return shardingIndex, shardingIndex > -1

}

func getSharding(index int, count int) int {
	if count <= 0 && index >= 0 {
		return index
	}
	if index < 0 {
		return -1
	}
	return index % count
}
