package model

/*
/api/friend/list get
{
    "code":200,
    "message":"成功",
    "data":[
        {
            "id":67,
            "nickName":"ironhead",
            "headImage":""
        }
    ]
}
*/

type ContactsInfo struct {
	Id        int    `json:"id"`
	NickName  string `json:"nickName"`
	HeadImage string `json:"headImage"`
}

type ContactsAdd struct {
	FriendId int `json:"friendId" form:"friendId"`
}
