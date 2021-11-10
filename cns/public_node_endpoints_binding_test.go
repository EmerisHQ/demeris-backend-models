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
			"PublicNodeEndpoints struct - Both empty",
			BothEndpointsEmpty,
			false,
		},
		{
			"PublicNodeEndpoints struct - Both filled",
			BothEndpointsFilled,
			false,
		},
		{
			"PublicNodeEndpoints struct - RPC empty",
			TendermintRPCEmpty,
			true,
		},
		{
			"PublicNodeEndpoints struct - API empty",
			CosmosAPIEmpty,
			true,
		},
	}
	// arrange
	validation.CosmosRPCURL(binding.Validator)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
	BothEndpointsEmpty = `
{
	"field1": "test",
	"public_node_endpoints": {}
}
`
	BothEndpointsFilled = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": "https://localhost:1234",
		"cosmos_api": "https://localhost:34567"
	}
}
`
	TendermintRPCEmpty = `
{
	"field1": "test",
	"public_node_endpoints": {
		"cosmos_api": "https://localhost:34567"
	}
}
`
	CosmosAPIEmpty = `
{
	"field1": "test",
	"public_node_endpoints": {
		"tendermint_rpc": "https://localhost:1234"
	}
}
`
)
