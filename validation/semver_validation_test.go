package validation_test

import (
	"testing"

	"github.com/emerishq/demeris-backend-models/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
)

func TestSemverValidation(t *testing.T) {
	type testStruct struct {
		Version string `binding:"semver"` //nolint: govet
	}

	tests := []struct {
		name      string
		version   string
		assertion require.ErrorAssertionFunc
	}{
		{
			"Cosmos SDK version string passes validation",
			"v0.44.3",
			require.NoError,
		},
		{
			"missing \"v\" before version string makes validation fail",
			"0.44.3",
			require.Error,
		},
		{
			"random string doesn't pass validation",
			"random string",
			require.Error,
		},
	}
	validation.Semver(binding.Validator)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := binding.Validator.ValidateStruct(testStruct{
				Version: tt.version,
			})

			tt.assertion(t, e)
		})
	}
}
