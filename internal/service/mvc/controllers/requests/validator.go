package requests

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func init() {
	_ = validate.RegisterValidation("bcrypt", bcryptValidator)
}

func bcryptValidator(fl validator.FieldLevel) bool {
	hash := fl.Field().String()

	_, err := bcrypt.Cost([]byte(hash))

	return err == nil
}

func BadRequest(err error) []*jsonapi.ErrorObject {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		return toJsonapiErrors(validationErrors)
	}

	return []*jsonapi.ErrorObject{
		{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Detail: err.Error(),
		},
	}
}

func toJsonapiErrors(m validator.ValidationErrors) []*jsonapi.ErrorObject {
	errs := make([]*jsonapi.ErrorObject, 0, len(m))

	for key, value := range m {
		errs = append(errs, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Meta: &map[string]interface{}{
				"field": key,
				"error": value.Error(),
			},
		})
	}

	return errs
}
