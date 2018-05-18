package controllers

import (
	"github.com/astaxie/beego"
	"loveHome/models"
)

type SessionController struct {
	beego.Controller
}

//将封装好的返回结构 变成json返回给前段
func (this *SessionController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *SessionController) DelSessionName() {
	beego.Info("========== /api/v1.0/session del Session succ ======")

	resp := make(map[string]interface{})

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//将session 删除
	this.DelSession("name")
	this.DelSession("user_id")
	this.DelSession("mobile")

	return
}

//房屋首页先不管，就返回成功
// /api/v1.0/session
func (this *SessionController) GetSessionName() {
	beego.Info("========== /api/v1.0/session get Session succ ======")

	resp := make(map[string]interface{})

	resp["errno"] = models.RECODE_SESSIONERR
	resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)

	defer this.RetData(resp)

	name_map := make(map[string]interface{})
	//将登陆或者注册给当前session存好的name字段返回给前段
	name := this.GetSession("name")
	if name != nil {
		resp["errno"] = models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		name_map["name"] = name.(string)
		resp["data"] = name_map
	}
	return
}
