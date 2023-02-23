package model

type GroupInfo struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	OwnerId        int    `json:"ownerId"`
	HeadImage      string `json:"headImage"`
	HeadImageThumb string `json:"headImageThumb"`
	Notice         string `json:"notice"`

	AliasName string `json:"aliasName"` // 用户在群的别名
	Remark    string `json:"remark"`    // 用户给这个群的别称
}

/*
创建群聊
http://127.0.0.1:8888/api/group/create?groupName=group1 POST

	{
		"code": 200,
		"message": "成功",
		"data": {
			"id": 26,
			"name": "grouptest3",
			"ownerId": 68,
			"headImage": "",
			"headImageThumb": "",
			"notice": null,
			"aliasName": "ironhead2",
			"remark": "grouptest3"
		}
	}
*/
type GroupCreateInput struct {
	GroupName string `json:"groupName" form:"groupName"`
}

/*
查询具体群聊

http://localhost/api/group/members/21 get
Content-Type: application/json;charset=UTF-8

	{
		"code": 200,
		"message": "成功",
		"data": [{
			"userId": 67,
			"aliasName": "ironhead",
			"headImage": "",
			"quit": false,
			"remark": "grouptest"
		}, {
			"userId": 68,
			"aliasName": "ironhead2",
			"headImage": "",
			"quit": false,
			"remark": "grouptest2"
		}]
	}
*/
type GroupMemberInfo struct {
	UserId    int    `json:"userId" form:"userId"`       //群成员的uid
	AliasName string `json:"aliasName" form:"aliasName"` //群成员在这个群的昵称
	HeadImage string `json:"headImage" form:"headImage"` //群成员的头像
	Quit      bool   `json:"quit" form:"quit"`           //是否退出群聊了
	Remark    string `json:"remark" form:"remark"`       //群成员在群的备注信息
}

/*
邀请群成员
{"groupId":21,"friendIds":[68]}
*/
type GroupMemberInviteInput struct {
	GroupId   int   `json:"groupId"`
	FriendIds []int `json:"friendIds"`
}
