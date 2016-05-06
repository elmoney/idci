package validators

type ParameterValidater interface {

	Validate(date []interface{}) error
}
