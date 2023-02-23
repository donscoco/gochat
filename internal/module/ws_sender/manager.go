package ws_sender

import concurrent_map "github.com/donscoco/gochat/pkg/container/concurrent-map"

// 1。创建一个concurrentmap 来管理 user_id 和 conn 之间的映射

var ConnChanMap *concurrent_map.ConcurrentMap // 因为im频繁连接和删除，所以使用concurrentmap

func InitConnChanMap() error {
	ConnChanMap = concurrent_map.New(8, nil) // todo 先用8个分区和默认hash，后续改成配置文件
	return nil
}

// 线程安全的
func SetChan(uid string, val chan []byte) {
	ConnChanMap.Set(uid, val)
}

func GetChan(uid string) (val chan []byte, ok bool) {
	data, ok := ConnChanMap.Get(uid)
	if !ok {
		return nil, ok
	}
	return data.(chan []byte), ok
}
