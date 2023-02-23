package model

/*

{"code":200,"message":"成功","data":{
	"id":67,"userName":"ironhead",
	"nickName":"ironhead","sex":0,"signature":"",
	"headImage":"","headImageThumb":"",
	"online":null
}}


{
    "userName":"ironhead",
    "nickName":"铁头班小朋友",
    "password":"123456",
    "confirmPassword":"123456"
}
{"code":200,"message":"成功","data":null}
*/

type UserInfo struct {
	Id             int    `json:"id"`
	UserName       string `json:"userName"`
	NickName       string `json:"nickName"`
	Sex            int    `json:"sex"`
	Signature      string `json:"signature"`
	HeadImage      string `json:"headImage"`
	HeadImageThumb string `json:"headImageThumb"`
	Online         bool   `json:"online"`
}

/*
http://localhost/api/user/findByNickName?nickName= get
{
    "code":200,
    "message":"成功",
    "data":[
        {
            "id":1,
            "userName":"blue",
            "nickName":"blue(作者)",
            "sex":0,
            "signature":"1",
            "headImage":"http://localhost/file/box-im/image/20221105/1667661524139.jpg",
            "headImageThumb":"http://localhost/file/box-im/image/20221105/1667661524185.jpg",
            "online":false
        },
        {
            "id":2,
            "userName":"张三",
            "nickName":"张三",
            "sex":1,
            "signature":"fadfa",
            "headImage":"http://localhost/file/box-im/image/20221229/1672284215704.jpg",
            "headImageThumb":"http://localhost/file/box-im/image/20221229/1672284215751.jpg",
            "online":true
        }
    ]
}
*/

type FindUserInput struct {
	NickName string `json:"nickName"  form:"nickName"` //form在 gin 解析的时候可以直接通过url中的参数解析
}
