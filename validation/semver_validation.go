package validation

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"golang.org/x/mod/semver"

	"github.com/go-playground/validator/v10"
)

const semverStr = "semver"

var ErrNotSemver = fmt.Errorf("version string is not semver-compliant")

func Semver(structValidator binding.StructValidator) {
	if v, ok := structValidator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation(semverStr, registerSemverValidation); err != nil {
			panic(err)
		}
	}
}

func registerSemverValidation(fl validator.FieldLevel) bool {
	path, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if err := validateSemver(path); err != nil {
		return false
	}

	return true
}

func validateSemver(value string) error {
	if !semver.IsValid(value) {
		return fmt.Errorf("version string is not semver-compliant")
	}

	return nil
}
