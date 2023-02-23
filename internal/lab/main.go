package main

import (
	"fmt"
	"github.com/donscoco/gochat/pkg/snowflake"
	"time"
)

// Redis的Sorted Set的分数范围从-(2^53)到+(2^53)。或者说是-9007199254740992 到 9007199254740992。更大的整数在内部用指数表示。

var aa = 4194304 // 10 位节点，12位id
var aaa int64 = int64(aa * 1000)
var aaaa = int64(aaa * 3600 * 60) // 缓存60天的数据，如果数据条数超出，再去查mysql

// 前 41 位是时间（毫秒单位）， 10位节点， 12 位id(1毫秒内的)
func main() {
	for i := 0; i < 100000; i++ {
		time.Sleep(1 * time.Second)
		sid := snowflake.GenID()

		// 方法1，切掉前面的时间
		fmt.Println(sid, " ", sid/aaaa, " ", sid%aaaa, " ", sid/aaa)
		// 方法2，切掉后面到的时间
		fmt.Println(sid, " ", sid/aaaa, " ", sid%aaaa, " ", sid/aaa)

	}

}

/*

 集合中最大的成员数为 232 - 1 (4294967295
不能直接存雪花算法直接生成的id，两个方法：
1。我们这里只是缓存，要的只是一段时间内的顺序性，可以% 保存60天内的数据
2。	或者，我们的score只是用来保存消息的时间顺序，不用特别严格，
	可以保存秒的时间单位。同一个毫秒下可能有多条，但是用户应该不关心一毫秒内哪个先。
	如果实在介意，可以在获取到数据后，检查每个时间点，对有相同分数的再进去value 查看 完整的雪花id进行排序。

原本的字段                    %后的保存的字段     按秒来存放

535721595272957952   591   293523848957952   127725981
535721599467261952   591   293528043261952   127725982
535721603686731776   591   293532262731776   127725983
535721607893618688   591   293536469618688   127725984
535721612092116992   591   293540668116992   127725985
535721616294809600   591   293544870809600   127725986
535721620489113600   591   293549065113600   127725987


方法2


*/