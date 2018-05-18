package controllers

import (
	"github.com/astaxie/beego"
	"loveHome/models"
)

type HousesIndexController struct {
	beego.Controller
}

//将封装好的返回结构 变成json返回给前段
func (this *HousesIndexController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//房屋首页先不管，就返回成功
func (this *HousesIndexController) GetHousesIndex() {
	beego.Info("========== /api/v1.0/houses/index  succ ======")

	resp := make(map[string]interface{})

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)
}
