package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const MESSAGE_BAD_REQUEST = "The server could not understand the request due to invalid syntax or missing parameters."
const MESSAGE_DATABASE_ERROR = "Failed to query database record. Contact the administrator for more information."
const MESSAGE_MISSING_PARAMS = "Missing one or more required parameter(s)"

// Returns string of missing required query parameters.
func MissingQueryParamsMessage(parameters []string) string {
	output := "Missing one or more required query parameter(s): "

	for i, v := range parameters {
		if i == len(parameters)-1 {
			output = output + v
		} else {
			output = output + v + ", "
		}
	}

	return output
}

// Returns string of invalid required query parameters.
func InvalidParamsMessage(fields []string) string {
	output := "One or more field(s) is of the wrong type or format: "

	for i, v := range fields {
		if i == len(fields)-1 {
			output = output + v
		} else {
			output = output + v + ", "
		}
	}

	return output
}

// Returns map of fields and "required" string.
func MissingParamsMessage(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, v := range verr {
		err := v.ActualTag()
		if v.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, v.Param())
		}
		errs[strings.ToLower(v.Field())] = err
	}

	return errs
}

// Returns map of field and its required data type.
func TypeErrorParamsMessage(jsErr json.UnmarshalTypeError) map[string]string {
	errs := make(map[string]string)
	errs[jsErr.Field] = jsErr.Type.String()

	return errs
}

/*
Returns map of MissingParamsMessage or TypeErrorParamsMessage.
If both error are not present, then return a generic error message.
*/
func GetErrorMessage(err error) (gin.H, string) {
	var validationError validator.ValidationErrors
	var unmarshalTypeError *json.UnmarshalTypeError

	if errors.As(err, &validationError) {
		return gin.H{MESSAGE_MISSING_PARAMS: MissingParamsMessage(validationError)}, ""
	} else if errors.As(err, &unmarshalTypeError) {
		return gin.H{MESSAGE_MISSING_PARAMS: TypeErrorParamsMessage(*unmarshalTypeError)}, ""
	} else {
		return nil, "One or more required fields are missing or invalid."
	}
}

// Return string of missing required pairs of fields.
func MissingValidPairMessage(field1 string, field2 string) string {
	return "Both fields " + field1 + " and " + field2 + " must be present and valid"
}
