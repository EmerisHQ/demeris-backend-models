package validation_test

import (
	"testing"

	"github.com/allinbits/demeris-backend-models/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
)

type test struct {
	Url string `binding:"cosmosrpcurl"`
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
			"Incorrect URL - HTTP",
			test{
				"http://123.123.123.123:1234",
			},
			"unsupported URL scheme http",
		},
		{
			"Incorrect URL - Other",
			test{
				"tcp://hostname:1234",
			},
			"unsupported URL scheme tcp",
		},
		{
			"Incorrect URL - Contains user",
			test{
				"https://user@hostname:1234",
			},
			"URL cannot contain user",
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
			"Incorrect URL - Pat",
			test{
				"https://hostname:123/test",
			},
			"URL cannot contain path info",
		},
		{
			"Incorrect URL - Fragment",
			test{
				"https://hostname:123#test",
			},
			"URL cannot contain path info",
		},
		{
			"Incorrect URL - Query",
			test{
				"https://hostname:123?test=1",
			},
			"URL cannot contain path info",
		},
	}
	// arrange
	validation.CosmosRPCURL(binding.Validator)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
