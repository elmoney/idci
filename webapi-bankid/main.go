package main

import (
	"github.com/idci/webapi-bankid/controllers"
	"github.com/idci/webapi-bankid/models"
	"gopkg.in/natefinch/lumberjack.v2"
        "log"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"

	"github.com/idci/webapi-bankid/helpers"
)

func main() {

	config:= models.InitConfig()


	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Settings.LogPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	log.Println("github.com/idci api init. Listening on :3000")


	_,err:= helpers.InitLedger(config)

	if(err!=nil) {
		log.Println(err)
		return
	}


	m:=martini.Classic()
	m.Map(config)
	m.Use(render.Renderer())

	m.Get("/", controllers.HomeIndex)

	// ** Units **
	m.Get("/api/v1/units", controllers.UnitsGet)
	m.Get("/api/v1/units/:id", controllers.UnitGet)
        m.Delete("/api/v1/units/:id", controllers.UnitRemove)
	m.Post("/api/v1/units",binding.Json(models.Unit{}), controllers.UnitCreate)
	m.Put("/api/v1/units",binding.Json(models.Unit{}), controllers.UnitUpdate)

	// ** Partners ***
	m.Post("/api/v1/partners", controllers.PartnerCreate)
	m.Delete("/api/v1/partners", controllers.PartnerRemove)

	// ** Requets **
	m.Post("/api/v1/request",binding.Json(models.CreateRequest{}), controllers.RequestCreate)
	m.Put("/api/v1/requestApprove",binding.Json(models.ApproveRequest{}), controllers.RequestApprove)
	m.Put("/api/v1/requestReject",binding.Json(models.RejectRequest{}), controllers.RequestReject)
	m.Get("/api/v1/request/:id", controllers.RequestGet)
	m.Get("/api/v1/requestsToApprove", controllers.RequestsToApprove)

	m.Run()
}


