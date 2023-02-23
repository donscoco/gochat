package data_manager

import (
	"errors"
	"golang.org/x/sync/singleflight"
)

var (
	ErrNoData = errors.New("空值")
)

// todo sigleflight 使用的id去lock，对于 userid=1和groupid=1 同样会锁
// todo 后续考虑是否使用一个singleflight，然后 key 加上前缀来区分比较好？

var _singleFlight *singleflight.Group

// var _singleFlightContacts *singleflight.Group
// var _singleFlightGroup *singleflight.Group
//var _singleFlightGroupMember *singleflight.Group

func init() {
	_singleFlight = new(singleflight.Group)
	//_singleFlightContacts = new(singleflight.Group)
	//_singleFlightGroup = new(singleflight.Group)
	//_singleFlightGroupMember = new(singleflight.Group)

}
