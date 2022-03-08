package validation_test

import (
	"testing"

	"github.com/emerishq/demeris-backend-models/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
)

type test struct {
	Url string `binding:"cosmosrpcurl"` //nolint: govet
}

func TestCosmosRpcValidation(t *testing.T) {
	tests := []struct {
		name        string
		testStruct  test
		expectedMsg string
	}{
		{
			"Correct URL - Local hostname",
			test{
				"https://hostname:1234",
			},
			"",
		},
		{
			"Correct URL - Full hostname",
			test{
				"https://hostname.url.com:1234",
			},
			"",
		},
		{
			"Correct URL - Uppercase protocol",
			test{
				"HTTPS://hostname.url.com:1234",
			},
			"",
		},
		{
			"Correct URL - IP",
			test{
				"https://123.123.123.123:1234",
			},
			"",
		},
		{
			"Correct URL - http",
			test{
				"http://123.123.123.123:1234",
			},
			"",
		},
		{
			"Correct URL - Contains user with pwd",
			test{
				"https://user:password@hostname:1234",
			},
			"",
		},
		{
			"Correct URL - Contains path",
			test{
				"https://hostname.com:1234/some/path",
			},
			"",
		},
		{
			"Correct URL - HTTP, path, user",
			test{
				"http://user:pwd@hostname.com:1234/some/path",
			},
			"",
		},
		{
			"Accepted value - Empty",
			test{
				"",
			},
			"",
		},
		{
			"Incorrect URL - Malformed",
			test{
				"hostname:123",
			},
			"malformed url",
		},
		{
			"Incorrect URL - Other",
			test{
				"tcp://hostname:1234",
			},
			"unsupported URL scheme tcp",
		},
		{
			"Incorrect URL - Contains user without pwd",
			test{
				"https://user@hostname:1234",
			},
			"URL cannot contain user without password",
		},
		{
			"Incorrect URL - No port",
			test{
				"https://hostname",
			},
			"invalid port ",
		},
		{
			"Incorrect URL - Invalid port",
			test{
				"https://hostname:foo",
			},
			"invalid port foo",
		},
		{
			"Incorrect URL - Port number too large",
			test{
				"https://hostname:70000",
			},
			"invalid port 70000",
		},
		{
			"Incorrect URL - Fragment",
			test{
				"https://hostname:123/path#test",
			},
			"URL cannot contain fragments or query params",
		},
		{
			"Incorrect URL - Query",
			test{
				"https://hostname:123/path?test=1",
			},
			"URL cannot contain fragments or query params",
		},
	}
	// arrange
	validation.CosmosRPCURL(binding.Validator)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// act
			e := binding.Validator.ValidateStruct(tt.testStruct)

			// assert
			if tt.expectedMsg == "" && e != nil {
				require.FailNow(t, "Unexpected test failure")
			} else if tt.expectedMsg != "" {
				// TODO: Properly assert expected error msg
				require.NotNil(t, e)
			}
		})
	}
}
