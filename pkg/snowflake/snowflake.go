package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node
var startTime = "2019-01-01" // 直接写死

func init() { // todo mark 删除，要让每台机器自己初始化 存放不同的machine id
	/*
	   因为到时候是部署在pod上的，每个machine id 无法区分，可以在redis上存放一个整数，每次机器获取就incr，让、然后%pod数
	*/

	if err := Init(startTime, 1); err != nil {
		fmt.Println("Init() failed, err = ", err)
		panic(err)
	}

}

// 开始时间，一定要是全局统一的一个时间点
// 机器节点，每个节点自己选择一个数字，保证每个节点的数字不一样就行
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

// 生成 64 位的 雪花 ID
func GenID() int64 {
	return node.Generate().Int64()
}
