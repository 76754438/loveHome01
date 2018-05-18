package controllers

import (
	"encoding/json" //添加这行
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"path"
)

type UserController struct {
	beego.Controller
}

//将封装好的返回结构 变成json返回给前段
func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *UserController) Reg() {
	beego.Info("==========/api/v1.0/users post succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)
	//定义一个make存储前端的信息
	var regRequestMap = make(map[string]interface{})

	//1 得到客户端请求的json数据 post数据
	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestMap)

	beego.Info("mobile = ", regRequestMap["mobile"])
	beego.Info("password = ", regRequestMap["password"])
	beego.Info("sms_code = ", regRequestMap["sms_code"])

	//2 判断数据的合法性
	if regRequestMap["mobile"] == "" || regRequestMap["password"] == "" || regRequestMap["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}
	//3 将数据存入到mysql数据库 user表中
	//现在只需要用户名和密码，user表中的其他信息由后面的业务完成
	user := models.User{}
	user.Mobile = regRequestMap["mobile"].(string) //注意表中是string类型
	//应该将password进行md5，SHA246,SHA1
	user.Password_hash = regRequestMap["password"].(string)
	user.Name = regRequestMap["mobile"].(string)

	//操作入数据库
	o := orm.NewOrm()

	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("insert error = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	//显示注册成功
	beego.Info("reg succ !!! user id = ", id)

	//4 将当前的用户的信息存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)
	return
}

//处理上传头像的业务
func (this *UserController) UploadAvatar() {
	//先写1-7的中文伪代码
	//7.返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//1.得到文件二进制数据,先看开发文档再写代码
	//进入bee.me看快速开发中的请求数据处理->文件上传
	//看图得到二进制数据图1-2，       avatar是前端制定的表单中的名字
	file, header, err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	//写这个之前去写06fafs_client.go
	fileBuffer := make([]byte, header.Size)
	if _, err := file.Read(fileBuffer); err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		return
	}
	//获取后缀，因为下面的FDFSUploadByBuffer函数中要用
	suffix := path.Ext(header.Filename) //把这个个字符串home.jpg.rmvb变成 .rmvb

	//2.将文件的二进制数据上传到fastdfs中 ---> fileid
	//fileBuffer--->fastdfs  ====>fileid                     suffix[1:]表示去掉.
	groupName, fileId, err := models.FDFSUploadByBuffer(fileBuffer, suffix[1:]) //"rmvb"
	if err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		beego.Info("upload file to fastdfs error err = ", err)
		return
	}
	//打印一下获取到的groupName和fileId
	beego.Info("fdfs upload succ groupname = ", groupName, "  fileid = ", fileId)

	//3.fileid --储存到--> user 表里avatar_ur字段中，
	//进入数据库看avatar_ur字段:mysql -uroot 123456
	//use lovehome   desc user
	//可以从seession中获得user.Id
	user_id := this.GetSession("user_id")                      //这里返回的是interface类型，所以下面的要转换
	user := models.User{Id: user_id.(int), Avatar_url: fileId} //更新Avatar_url字段
	//下面的操作是更新Avatar_url字段

	//4.数据库的操作，
	o := orm.NewOrm()
	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//5.将fileid拼接成一个完整的url路径，这个ip先写死
	avatar_url := "http://192.168.146.155:8080/" + fileId

	//6.安装协议做出json返回给前端
	//看讲义中的返回成功返回的数据格式，发现要map，
	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url
	resp["data"] = url_map

	return

}
func (this *UserController) Login() {
	beego.Info("==========/api/v1.0/sessions login succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	var loginRequestMap = make(map[string]interface{})

	//1 得到客户端请求的json数据 post数据
	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap)

	beego.Info("mobile = ", loginRequestMap["mobile"])
	beego.Info("password = ", loginRequestMap["password"])

	//2 判断数据的合法性
	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//3 查询数据库得到user
	var user models.User

	o := orm.NewOrm()
	//select password from user where user.name = name
	qs := o.QueryTable("user")
	if err := qs.Filter("mobile", loginRequestMap["mobile"]).One(&user); err != nil {
		//查询失败
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}

	//4 对比密码
	if user.Password_hash != loginRequestMap["password"].(string) {
		resp["errno"] = models.RECODE_PWDERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PWDERR)
		return
	}

	beego.Info("==== login succ!!! === user.name = ", user.Name)

	//5 将当前的用户的信息存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}
