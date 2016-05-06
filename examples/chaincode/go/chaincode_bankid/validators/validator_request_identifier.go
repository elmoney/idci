package validators

import (
	"github.com/op/go-logging"
	"regexp"
	"errors"
)

type ValidatorRequestIdentifier struct {

	Log *logging.Logger
}

func (v *ValidatorRequestIdentifier) Validate(data []interface{}) error {

	requestId := data[0].(string)
	methodName := data[1].(string)

	requestIdValid, err := regexp.MatchString("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}",
		requestId)

	if err != nil || !requestIdValid {
		if v.Log != nil {
			v.Log.Info("[%s]RequestId contain incorrect format", methodName)
		}

		return errors.New("RequestId contain incorrect format")
	}

	return nil
}