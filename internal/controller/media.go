package controller

import (
	"github.com/donscoco/gochat/internal/base/oss"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

type MediaController struct{}

// https://gin-gonic.com/zh-cn/docs/examples/upload-file/single-file/
func MediaRegister(group *gin.RouterGroup) {
	mediaController := &MediaController{}
	group.POST("/image/upload", mediaController.UploadImage) // 上传图片
	group.POST("/file/upload", mediaController.UploadFile)   // 上传文件
}

/*
/api/image/upload
Request Method: POST
Content-Type: multipart/form-data; boundary=----WebKitFormBoundarympYRd449IlAxLtTS

{"code":200,"message":"成功","data":{"originUrl":"http://localhost/file/box-im/image/20230203/1675414276279.png","thumbUrl":"http://localhost/file/box-im/image/20230203/1675414276293.png"}}
*/
func (mc *MediaController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}
	f, err := file.Open()
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}

	err = oss.PutObject(oss.DefaultOSS, oss.DefaultBucket, "user_img/"+file.Filename, f, time.Now(), true)
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}

	out := model.UploadFileOutput{
		OriginUrl: oss.DefaultURLPrefix + "user_img/" + file.Filename,
		ThumbUrl:  oss.DefaultURLPrefix + "user_img/" + file.Filename,
	}

	bl.ResponseSuccess(c, out)

}

/*
Request URL: /api/file/upload
Request Method: POST

{"code":200,"message":"",
"data":"http://localhost/file/box-im/file/20230204/1675522180846.wav"}
*/

func (mc *MediaController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}
	f, err := file.Open()
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}

	err = oss.PutObject(oss.DefaultOSS, oss.DefaultBucket, "user_file/"+file.Filename, f, time.Now(), true)
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}

	//out := model.UploadFileOutput{
	//	OriginUrl: oss.DefaultURLPrefix + "user_file/" + file.Filename,
	//	ThumbUrl:  oss.DefaultURLPrefix + "user_file/" + file.Filename,
	//}

	bl.ResponseSuccess(c, oss.DefaultURLPrefix+"user_file/"+file.Filename)
}
