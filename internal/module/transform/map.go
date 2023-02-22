package transform

import "github.com/donscoco/gochat/internal/dao"

func GroupMemberGroupByGroupIdListToMap(member []dao.GroupMember) (m map[int]*dao.GroupMember) {
	m = make(map[int]*dao.GroupMember, 16)

	for _, e := range member {
		m[e.GroupId] = &e
	}
	return
}
