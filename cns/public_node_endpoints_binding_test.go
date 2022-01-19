package cns_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/allinbits/demeris-backend-models/validation"

	"github.com/gin-gonic/gin/binding"
)

// Validate that both fields are filled or empty
func TestPublicNodeEndpointsBinding(t *testing.T) {

	tests := []struct {
		name     string
		testJson string
		fails    bool
	}{
		{
			"PublicNodeEndpoints struct - Missing",
			CompletelyMissing,
			false,
		},
		{
			"PublicNodeEndpoints struct - Both not defined",
			BothEndpointsNotDefined,
			false,
		},
		{
			"PublicNodeEndpoints struct - Both empty",
			BothEndpointsEmpty,
			true,
		},
		{
			"PublicNodeEndpoints struct - Both filled",
			BothEndpointsFilled,
			false,
		},
		{
			"PublicNodeEndpoints struct - RPC not defined",
			TendermintRPCNotDefined,
			true,
		},
		{
			"PublicNodeEndpoints struct - RPC empty",
			TendermintRPCEmpty,
			true,
		},
		{
			"PublicNodeEndpoints struct - API not defined",
			CosmosAPINotDefined,
			true,
		},
		{
			"PublicNodeEndpoints struct - API empty",
			CosmosAPIEmpty,
			true,
		},
		{
			"PublicNodeEndpoints struct - Multiple RPCs",
			TendermintMultipleRPC,
			false,
		},
		{
			"PublicNodeEndpoints struct - Multiple APIs",
			CosmosMultipleAPI,
			false,
		},
		{
			"PublicNodeEndpoints struct - Invalid RPC value",
			TendermintRPCInvalid,
			true,
		},
		{
			"PublicNodeEndpoints struct - Invalid API value",
			CosmosAPIInvalid,
			true,
		},
	}
	// arrange
	validation.CosmosRPCURL(binding.Validator)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct := struct {
				Field1              string                  `json:"field1"`
				PublicNodeEndpoints cns.PublicNodeEndpoints `binding:"dive" json:"public_node_endpoints,omitempty"`
			}{}

			req := requestWithBody(tt.testJson)
			err := binding.JSON.Bind(req, &testStruct)

			// assert
			if tt.fails {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

const (
	CompletelyMissing = `
{
	"field1": "test"
}
`

	BothEndpointsNotDefined = `
{
	"field1": "test",
	"public_node_endpoints": {}
}
`

	BothEndpointsEmpty = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": [],
		"cosmos_api": []
	}
}
`
	BothEndpointsFilled = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https://localhost:1234"],
		"cosmos_api": ["https://localhost:34567"]
	}
}
`
	TendermintRPCNotDefined = `
{
	"field1": "test",
	"public_node_endpoints": {
		"cosmos_api": ["https://localhost:34567"]
	}
}
`
	TendermintRPCEmpty = `
{
	"field1": "test",
	"public_node_endpoints": {
		"cosmos_api": ["https://localhost:34567"],
		"tendermint_rpc": []
	}
}
`
	CosmosAPINotDefined = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https://localhost:1234"]
	}
}
`
	CosmosAPIEmpty = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https://localhost:1234"],
		"cosmos_api": []
	}
}
`
	TendermintMultipleRPC = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https://localhost:1234","https://localhost:1235"],
		"cosmos_api": ["https://localhost:34567"]
	}
}
`
	CosmosMultipleAPI = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https://localhost:1234"],
		"cosmos_api": ["https://localhost:34567","https://localhost:34566"]
	}
}
`
	TendermintRPCInvalid = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https"],
		"cosmos_api": ["https://localhost:34567"]
	}
}
`
	CosmosAPIInvalid = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": ["https://localhost:1234"],
		"cosmos_api": ["https"]
	}
}
`
)
