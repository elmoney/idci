package helpers

import (
	"github.com/idci/webapi-bankid/utils"
	"encoding/json"
	"strconv"
	"errors"
	"net/http"
	"log"
	"fmt"
	"github.com/idci/webapi-bankid/models"
)


//UnitList список участников
func UnitList(config models.Config) (*[]models.Unit, error) {
	status, resp := utils.QueryStateGet(config.Settings.ChainCodeName, config.Settings.PeerQueryURL, "state-view", "_uContainerIdx")

	log.Println("UnitRemove")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		var request models.UnitsGetRequest

		err := json.Unmarshal([]byte(resp), &request)
		if (err != nil) {
			return nil, err
		}

		if (request.Error != "") {
			return nil, errors.New(request.Error)
		}

		return &request.OK.Units, nil

	}

	return nil, errors.New(strconv.Itoa(status))

}

//UnitRemove удаление
func UnitRemove(config models.Config, alias string, sign string) (error) {

	status, resp := utils.InvokeCC(`unit-remove`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName, alias, sign)

	log.Println("UnitList")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}
	return errors.New(strconv.Itoa(status))

}

//UnitCreate создание
func UnitCreate(config models.Config, unit models.Unit, sign string) (error) {

	status, resp := utils.InvokeCC(`unit-create`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unit.Alias, unit.FullName, strconv.Itoa(unit.Status), unit.Identifier.Name, unit.Identifier.Description, unit.Identifier.Regexp, fmt.Sprintf("%x", unit.PublicKey), sign)

	log.Println("UnitCreate")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}
	return errors.New(strconv.Itoa(status))

}

//UnitUpdate  обновление
func UnitUpdate(config models.Config, unit models.Unit, sign string) (error) {

	status, resp := utils.InvokeCC(`unit-update`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unit.Alias, unit.FullName, strconv.Itoa(unit.Status), unit.Identifier.Name, unit.Identifier.Description, unit.Identifier.Regexp, fmt.Sprintf("%x", unit.PublicKey), sign)

	log.Println("UnitUpdate")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}
	return errors.New(strconv.Itoa(status))

}



//PartnerCreate создание участника
func PartnerCreate(config models.Config, unit string, partner string, sign string) (error) {

	status, resp := utils.InvokeCC(`unit-partner-create`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unit, partner, sign)

	log.Println("PartnerAdd")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}
	return errors.New(strconv.Itoa(status))

}

//PartnerRemove удаление  партнера
func PartnerRemove(config models.Config, unit string, partner string, sign string) (error) {

	status, resp := utils.InvokeCC(`unit-partner-remove`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unit, partner, sign)

	log.Println("PartnerRemove")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}
	return errors.New(strconv.Itoa(status))

}
