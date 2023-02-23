package transform

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/conn_manager"
)

func TransformContactsList(srcList []dao.Contacts, usermap map[int]dao.User) (tgt []model.ContactsInfo) {
	tgt = make([]model.ContactsInfo, 0, 16)
	for _, c := range srcList {
		u := usermap[c.FriendId]
		tgt = append(tgt, model.ContactsInfo{
			Id:        c.FriendId,
			NickName:  c.FriendNickName,
			HeadImage: u.HeadImage,
		})
	}

	return
}

func TransformGroupList(srcList []dao.Group, m map[int]*dao.GroupMember) (tgt []model.GroupInfo) {
	tgt = make([]model.GroupInfo, 0, 16)
	for _, g := range srcList {
		member, ok := m[g.Id]
		if !ok {
			continue
		}

		tgt = append(tgt, model.GroupInfo{
			Id:             g.Id,
			Name:           g.GroupName,
			OwnerId:        g.OwnerId,
			HeadImage:      g.HeadImage,
			HeadImageThumb: g.HeadImage,
			Notice:         g.Notice,
			AliasName:      member.AliasName,
			Remark:         member.Remark,
		})
	}
	return
}

func TransformUserList(srcList []dao.User) (tgt []model.UserInfo) {
	tgt = make([]model.UserInfo, 0, 16)
	for _, u := range srcList {
		tgt = append(tgt, model.UserInfo{
			Id:             u.Id,
			UserName:       u.UserName,
			NickName:       u.NickName,
			Sex:            u.Sex,
			Signature:      u.Signature,
			HeadImage:      u.HeadImage,
			HeadImageThumb: u.HeadImage,
			Online:         conn_manager.IsAlive(u.Id),
		})
	}
	return
}

func TransformGroupMemberList(srcList []dao.GroupMember, usermap map[int]dao.User) (tgt []model.GroupMemberInfo) {
	tgt = make([]model.GroupMemberInfo, 0, 16)
	for _, gm := range srcList {
		u := usermap[gm.UserId]
		tgt = append(tgt, model.GroupMemberInfo{
			UserId:    gm.UserId,
			AliasName: gm.AliasName,
			HeadImage: u.HeadImage,
			Quit:      false, // 查询is delete 来获取是否退出群聊
			Remark:    gm.Remark,
		})
	}
	return
}
