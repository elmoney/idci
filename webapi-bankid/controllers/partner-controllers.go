package controllers

import (
	"net/http"
	"github.com/idci/webapi-bankid/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"fmt"
	"strings"
	"github.com/idci/oaep"
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"obc-bankid/webapi-bankid/Const"
	"github.com/idci/webapi-bankid/helpers"
)


//PartnerRemove удаление партнера
func PartnerRemove(r *http.Request, ren render.Render, config *models.Config, params martini.Params) {

	log.Println(r)
	log.Println(params)


	body, err := ioutil.ReadAll(r.Body)

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)
		return
	}

	v, err := url.ParseQuery(string(body))

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)
		return
	}

	unitAlias := v.Get("unitAlias")
	partnerAlias := v.Get("partnerAlias")

	log.Println(r)


	if(unitAlias=="" || partnerAlias==""){
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
		if val.Alias == unitAlias {
			exists =1
			break
		}
	}

	if(exists==0){
		ren.JSON(http.StatusNotFound, nil)
		return;
	}

	exists=0

	for _, val := range *units {
		if val.Alias == partnerAlias {
			exists =1
			break
		}
	}

	if(exists==0){
		ren.JSON(http.StatusNotFound, nil)
		return;
	}


	plain:= fmt.Sprintf("%s|%s", strings.ToLower(unitAlias),strings.ToLower(partnerAlias))
	signedData, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(plain));

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	//fmt.Println(signedData)
	hexed :=fmt.Sprintf("%x",signedData)
	//fmt.Println(hexed)
	log.Println(hex.DecodeString(hexed))

	err= helpers.PartnerRemove(*config,unitAlias,partnerAlias,hexed)

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	ren.JSON(http.StatusOK, nil)
	return

}



//PartnerCreate добавление партнера
func PartnerCreate(r *http.Request, ren render.Render, config *models.Config, params martini.Params) {
	log.Println(r)
	log.Println(params)


	body, err := ioutil.ReadAll(r.Body)

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)
		return
	}

	v, err := url.ParseQuery(string(body))

	if (err != nil) {
		log.Println(err)
		ren.JSON(http.StatusInternalServerError, err)
		return
	}

	unitAlias := v.Get("unitAlias")
	partnerAlias := v.Get("partnerAlias")

	log.Println(r)


	if(unitAlias=="" || partnerAlias==""){
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

	var unit models.Unit

	for _, val := range *units {
		if val.Alias == unitAlias {
			unit = val
			exists =1
			break
		}
	}

	if(exists==0){
		ren.JSON(http.StatusNotFound, nil)
		return;
	}

	exists=0

	for _, val := range *units {
		if val.Alias == partnerAlias {
			exists =1
			break
		}
	}

	if(exists==0){
		ren.JSON(http.StatusNotFound, nil)
		return;
	}


	for _, val := range unit.Partners {
		if strings.ToLower(val) == strings.ToLower(partnerAlias) {
			ren.JSON(http.StatusBadRequest, Errors.PartnerAlreadyExists)
			return;
		}
	}


	plain:= fmt.Sprintf("%s|%s", strings.ToLower(unitAlias),strings.ToLower(partnerAlias))
	signedData, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(plain));

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	//fmt.Println(signedData)
	hexed :=fmt.Sprintf("%x",signedData)
	//fmt.Println(hexed)
	log.Println(hex.DecodeString(hexed))

	err= helpers.PartnerCreate(*config,unitAlias,partnerAlias,hexed)

	if(err!=nil){
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	ren.JSON(http.StatusOK, nil)
	return
}
