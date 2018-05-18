package controllers

import (
	//	"fmt"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"         //添加这行，要用到
	_ "github.com/astaxie/beego/cache/redis" //添加这行，上面的那行需要用到这个包
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"time"
)

type AreaController struct {
	beego.Controller
}

func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//  /api/v1.0/areas [get]
func (this *AreaController) GetAreaInfo() {
	beego.Info("==========/api/v1.0/area get succ!!!=========")
	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//0  链接redis数据
	cache_conn, err := cache.NewCache("redis", `{"key":"lovehome","conn":"127.0.0.1:6380","dbNum":"0"}`)
	if err != nil {
		beego.Info("cache redis conn err, err= ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return

	}
	/*
		cache_conn.Put("xixi", "lala", time.Second*300)
		//1 从缓存中redis读数据
		values := cache_conn.Get("xixi")

		fmt.Println(values) //打印出来是空
		if values != nil {
			beego.Info(" cache get value = ", values)
			fmt.Printf("value = %s\n", values) //添加fmt包,打印出heihei
		}*/
	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value != nil {
		//2 如果redis有 之前的json字符串数据那么直接返回给前段
		//说明area_info key是存在的  value就是要返回给前段的json值
		beego.Info(" ====== get area_info from cache !!! ======")

		var area_info interface{}

		json.Unmarshal(areas_info_value.([]byte), &area_info)
		resp["data"] = area_info
		return
	}
	//2 如果redis有 之前的json字符串数据那么直接返回给前段

	//3 如果redis没有之前的json字符串数据， 从mysql查
	o := orm.NewOrm()

	//得到查到的areas数据，这数据格式是怎么样的，看models.go中的Area结构体，
	//Houses字段我们可以先不管，看数据库中的area表
	var areas []models.Area //[{aid,aname},{aid,aname},{aid,aname}]

	qs := o.QueryTable("area") //查数据库中的表
	num, err := qs.All(&areas) //把数据库中的表中的内容放到areas中
	if err != nil {
		//返回错误信息给前端,去写第二步，然后接着往下写
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)

		return
	}
	if num == 0 {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)

		return
	}

	//succ 把数组储存到map结构体的data中
	resp["data"] = areas
	//将 areas json字符串 存到area_info redis的key中
	areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info --> redis fail err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errno"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	//将封装好的返回结构体map 发送给前段

	return

}
