package controllers

import (
	"net/http"
	"github.com/idci/webapi-bankid/models"
	"log"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/idci/oaep"
	"fmt"
	"strings"

	"encoding/hex"
	"obc-bankid/webapi-bankid/Const"
	"github.com/idci/webapi-bankid/helpers"
)



//UnitsGet  получить всех участников системы
func UnitsGet(r *http.Request, ren render.Render, config *models.Config) {
	log.Println(r)
	units,err := helpers.UnitList(*config)

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)

	} else {
		log.Println(units)
		ren.JSON(http.StatusOK, units)
	}
}


//UnitGet получить учатника системы
func UnitGet(r *http.Request, ren render.Render, config *models.Config, params martini.Params) {
	log.Println(r)
	log.Println(params)
	 units,err := helpers.UnitList(*config)

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)

	} else {
		alias := params["id"]
		for _, val := range *units {
			if val.Alias == alias {
				ren.JSON(http.StatusOK, val)
				return;
			}
		}
		ren.JSON(http.StatusNotFound, nil)
	}
}


//UnitRemove удалить участника системы
func UnitRemove(r *http.Request, ren render.Render, config *models.Config, params martini.Params) {
	log.Println(r)
	log.Println(params)

	alias := params["id"]

	signedData, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(alias));

	hexed :=fmt.Sprintf("%x",signedData)
	log.Println( hex.DecodeString(hexed))

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
	}

	units,err := helpers.UnitList(*config)
	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)

	} else {
		alias := params["id"]

		//проверка не явлется ли данный юнит партнером
		for _, val := range *units {
			for _, partner := range val.Partners {
				if (partner == alias) {
					ren.JSON(http.StatusBadRequest, Errors.UnitIsPartner)
					return;
				}
			}
		}

		for _, val := range *units {
			if val.Alias == alias {
				helpers.UnitRemove(*config, alias, hexed)
				ren.JSON(http.StatusOK, nil)
			}
		}
		ren.JSON(http.StatusNotFound, nil)
	}

}


//UnitCreate сохдать участника
func UnitCreate(r *http.Request, ren render.Render, params martini.Params,
config *models.Config,
unit models.Unit) {
	log.Println(r)
	log.Println(unit.Alias)

	if(unit.Alias==""){
		ren.JSON(http.StatusBadRequest, Errors.UnitAliasIsEmpty)
		return;
	}

	 units,err := helpers.UnitList(*config)

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)
                 return
	}

	for _, val := range *units {
		if val.Alias == unit.Alias {
			ren.JSON(http.StatusBadRequest, Errors.UnitWithSameAliasAlreadyExists)
			return;
		}
	}

	plain:= fmt.Sprintf("%s|%s", strings.ToLower(unit.Alias),strings.ToLower(unit.FullName))
	signedData, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(plain));

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	//fmt.Println(signedData)
	hexed :=fmt.Sprintf("%x",signedData)
	//fmt.Println(hexed)
	log.Println( hex.DecodeString(hexed))

	unit.Identifier.Regexp = fmt.Sprintf("%x",[]byte(unit.Identifier.Regexp))

	log.Println(plain)
	log.Println(hexed)

       err= helpers.UnitCreate(*config,unit,hexed)

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	ren.JSON(http.StatusOK, nil)
	return

}



//UnitUpdate  обновить участника
func UnitUpdate(r *http.Request, ren render.Render, params martini.Params,
config *models.Config,
unit models.Unit) {

	log.Println(r)
	log.Println(unit.Alias)

	if(unit.Alias==""){
		ren.JSON(http.StatusBadRequest, Errors.UnitAliasIsEmpty)
		return;
	}

	 units,err := helpers.UnitList(*config)

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)
		return
	}

	exists:=0

	for _, val := range *units {
		if val.Alias == unit.Alias {
			exists =1
			break
		}
	}

	if(exists==0){
		ren.JSON(http.StatusNotFound, nil)
		return;
	}


	plain:= fmt.Sprintf("%s|%s", strings.ToLower(unit.Alias),strings.ToLower(unit.FullName))
	signedData, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(plain));

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	//fmt.Println(signedData)
	hexed :=fmt.Sprintf("%x",signedData)
	//fmt.Println(hexed)
	log.Println( hex.DecodeString(hexed))

	unit.Identifier.Regexp = fmt.Sprintf("%x",[]byte(unit.Identifier.Regexp))

	log.Println(plain)
	log.Println(hexed)

	err= helpers.UnitUpdate(*config,unit,hexed)

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	ren.JSON(http.StatusOK, nil)
	return

}

