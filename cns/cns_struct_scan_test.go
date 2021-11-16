package cns_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allinbits/demeris-backend-models/cns"
)

func TestPublicNodeEndpointScan(t *testing.T) {

	tests := []struct {
		name     string
		testJson string
		fails    bool
	}{
		{
			"PublicNodeEndpoints - Full JSON",
			PublicNodeEndpointsFullJSON,
			false,
		},
		{
			"Chain struct - Empty string",
			"",
			true,
		},
		{
			"Chain struct - Invalid string",
			"foo",
			true,
		},
	}
	// arrange

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pne := cns.PublicNodeEndpoints{}
			err := pne.Scan([]byte(tt.testJson))

			// assert
			if tt.fails {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, pne.CosmosAPI)
				assert.NotEmpty(t, pne.TendermintRPC)
			}
		})
	}
}

const (
	PublicNodeEndpointsFullJSON = `{
		"tendermint_rpc": "https://localhost:1234",
		"cosmos_api": "https://127.0.01:2345"
	}`
)
