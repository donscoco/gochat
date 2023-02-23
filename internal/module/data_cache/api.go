package data_cache

import (
	"errors"
	"time"
)

const (
	EMPTY = "nil"
	TTL   = 5 * 60 * time.Second

	PREFIX_USER_CACHE     = "c_user_"
	PREFIX_CONTACTS_CACHE = "c_contacts_"
	PREFIX_GROUP_CACHE    = "c_group_"

	PREFIX_GROUP_MEMBER_CACHE = "c_gm_" // c_gm_{uid}_{gid}

	PREFIX_CONVERSATION_CACHE = "c_conv_" // c_conv_{uid}_{gid},c_conv_{uid}_{fid} // 存放读写offset
	PREFIX_MSG_CACHE          = "c_msg_"  // c_msg_{uid}_{gid},c_msg_{uid}_{fid} 约定小的uid在前面

	PREFIX_CONTACTS_LIST_CACHE     = "c_contacts_list_" // c_gm_list_{uid} 用户的通讯录列表
	PREFIX_GROUP_MEMBER_LIST_CACHE = "c_gm_list_"       // c_gm_list_{gid} 群组的成员列表
	PREFIX_GROUP_LIST_CACHE        = "c_group_list_"    // c_group_list_{uid} 用户所属的群组列表
)

var (
	ErrNoData = errors.New("空值") // 应对缓存穿透 // todo 后续考虑是否改成布隆过滤器来应对缓存穿透
)
