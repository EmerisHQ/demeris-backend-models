package validation_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/allinbits/demeris-backend-models/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
)

type test struct {
	url string `binding:cosmosrpcurl` //nolint: govet
}

func TestCosmosRpcValidation(t *testing.T) {
	tests := []struct {
		name          string
		testStruct    test
		expectedError error
	}{
		{
			"Correct URL - Local hostname",
			test{
				"https://hostname:1234",
			},
			nil,
		},
		{
			"Correct URL - Full hostname",
			test{
				"https://hostname.url.com:1234",
			},
			nil,
		},
		{
			"Correct URL - Uppercase protocol",
			test{
				"HTTPS://hostname.url.com:1234",
			},
			nil,
		},
		{
			"Correct URL - IP",
			test{
				"https://123.123.123.123:1234",
			},
			nil,
		},
		{
			"Incorrect URL - Empty",
			test{
				"",
			},
			errors.New("malformed url"),
		},
		{
			"Incorrect URL - Malformed",
			test{
				"hostname:123",
			},
			errors.New("malformed url"),
		},
		{
			"Incorrect URL - HTTP",
			test{
				"http://123.123.123.123:1234",
			},
			errors.New("unsupported URL scheme http"),
		},
		{
			"Incorrect URL - Other",
			test{
				"tcp://hostname:1234",
			},
			fmt.Errorf("unsupported URL scheme tcp"),
		},
		{
			"Incorrect URL - Contains user",
			test{
				"https://user@hostname:1234",
			},
			errors.New("URL cannot contain user"),
		},
		{
			"Incorrect URL - No port",
			test{
				"https://hostname",
			},
			errors.New("invalid port "),
		},
		{
			"Incorrect URL - Invalid port",
			test{
				"https://hostname:foo",
			},
			errors.New("invalid port foo"),
		},
		{
			"Incorrect URL - Port number too large",
			test{
				"https://hostname:70000",
			},
			errors.New("invalid port 70000"),
		},
		{
			"Incorrect URL - Pat",
			test{
				"https://hostname:123/test",
			},
			errors.New("URL cannot contain path info"),
		},
		{
			"Incorrect URL - Fragment",
			test{
				"https://hostname:123#test",
			},
			errors.New("URL cannot contain path info"),
		},
		{
			"Incorrect URL - Query",
			test{
				"https://hostname:123?test=1",
			},
			errors.New("URL cannot contain path info"),
		},
	}
	validation.CosmosRPCURL(binding.Validator)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedError, binding.Validator.ValidateStruct(tt.testStruct))
		})
	}
}
