package validation_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin/binding"

	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/require"

	"github.com/allinbits/demeris-backend-models/validation"
)

// ----- Missing Fields test ----

type e struct {
	Arg string
}

func (e e) Tag() string {
	return "required"
}

func (e e) ActualTag() string {
	return "Actual" + e.Arg
}

func (e e) Namespace() string {
	panic("implement me")
}

func (e e) StructNamespace() string {
	panic("implement me")
}

func (e e) Field() string {
	return "Field" + e.Arg
}

func (e e) StructField() string {
	panic("implement me")
}

func (e e) Value() interface{} {
	panic("implement me")
}

func (e e) Param() string {
	panic("implement me")
}

func (e e) Kind() reflect.Kind {
	panic("implement me")
}

func (e e) Type() reflect.Type {
	panic("implement me")
}

func (e e) Translate(_ ut.Translator) string {
	panic("implement me")
}

func (e e) Error() string {
	panic("implement me")
}

func TestMissingFields(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		fieldName bool
		want      []string
	}{
		{
			"not validation error",
			fmt.Errorf("not validation"),
			false,
			nil,
		},
		{
			"validation error",
			validator.ValidationErrors{
				e{},
				e{Arg: "second"},
			},
			false,
			[]string{"Field", "Fieldsecond"},
		},
		{
			"validation error with field name",
			validator.ValidationErrors{
				e{},
				e{Arg: "second"},
			},
			true,
			[]string{"Actual", "Actualsecond"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, validation.MissingFields(tt.err, tt.fieldName))
		})
	}
}

func TestMissingFieldsErr(t *testing.T) {

	tests := []struct {
		name      string
		err       error
		fieldName bool
		want      error
	}{
		{
			"not validation error",
			fmt.Errorf("not validation"),
			false,
			fmt.Errorf("not validation"),
		},
		{
			"validation error",
			validator.ValidationErrors{
				e{},
				e{Arg: "second"},
			},
			false,
			fmt.Errorf("missing fields: Field,Fieldsecond"),
		},
		{
			"validation error with field name",
			validator.ValidationErrors{
				e{},
				e{Arg: "second"},
			},
			true,
			fmt.Errorf("missing fields: Actual,Actualsecond"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, validation.MissingFieldsErr(tt.err, tt.fieldName))
		})
	}
}

// --- DerivationPath test ---

type derivationTest struct {
	DerivationPath string `binding:"required,derivationpath"`
}

func TestDerivationPath(t *testing.T) {

	tests := []struct {
		name       string
		testStruct derivationTest
		fails      bool
	}{
		{
			"Valid - Ethereum",
			derivationTest{
				DerivationPath: "m/44'/60'/0'/1",
			},
			false,
		},
		{
			"Valid - Ethereum Classic",
			derivationTest{
				DerivationPath: "m/44'/61'/0'/0",
			},
			false,
		},
		{
			"Invalid - Mixed chars",
			derivationTest{
				DerivationPath: "m/abc'/61'/0'/0",
			},
			true,
		},
		{
			"Invalid - Empty",
			derivationTest{
				DerivationPath: "",
			},
			true,
		},
		{
			"Invalid - Borked",
			derivationTest{
				DerivationPath: "/////",
			},
			true,
		},
	}
	// arrange
	validation.DerivationPath(binding.Validator)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// act
			e := binding.Validator.ValidateStruct(tt.testStruct)

			// assert
			if !tt.fails && e != nil {
				require.FailNow(t, "Unexpected test failure")
			} else if tt.fails {
				require.NotNil(t, e)
			}
		})
	}
}
