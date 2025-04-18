package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const alphaSpaceRegexString string = "^[a-zA-Z ]+$"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("alphaspace", isAlphaSpace)

	return validate
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}

func ToErrResponse(err error) []string {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		resp := make([]string, len(fieldErrors))

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp[i] = fmt.Sprintf("%s is a required field", err.Field())
			case "max":
				resp[i] = fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param())
			case "url":
				resp[i] = fmt.Sprintf("%s must be a valid URL", err.Field())
			case "alphaspace":
				resp[i] = fmt.Sprintf("%s can only contain alphabetic and space characters", err.Field())
			case "datetime":
				if err.Param() == "2006-01-02" {
					resp[i] = fmt.Sprintf("%s must be a valid date", err.Field())
				} else {
					resp[i] = fmt.Sprintf("%s must follow %s format", err.Field(), err.Param())
				}
			default:
				resp[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}

		return resp
	}

	return nil
}
