package cns_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/allinbits/demeris-backend-models/validation"

	"github.com/gin-gonic/gin/binding"
)

func TestChainBinding(t *testing.T) {

	tests := []struct {
		name     string
		testJson string
		fails    bool
	}{
		{
			"Chain struct - Without PublicNodeEndpoints",
			ChainWithoutPublicNodeEndpoints,
			false,
		},
		{
			"Chain struct - With PublicNodeEndpoints",
			ChainWithPublicNodeEndpoints,
			false,
		},
	}
	// arrange
	// FIXME: add the other custom validators here after `utils` extraction
	validation.CosmosRPCURL(binding.Validator)
	validation.DerivationPath(binding.Validator)
	validation.JSONFields(binding.Validator)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnsStruct := cns.Chain{}

			req := requestWithBody(tt.testJson)
			err := binding.JSON.Bind(req, &cnsStruct)

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
	ChainWithoutPublicNodeEndpoints = `
{
	"enabled": true,
	"chain_name": "foo",
	"logo": "logo.png",
	"display_name": "FooBar",
	"primary_channel": {
		"key1": "value1",
		"key2": "value2"
	},
	"denoms": [
	{
		"name": "denom1",
		"display_name": "Denom 1",
		"logo": "https://logo.com",
		"precision": 12,
		"verified": true,
		"stakable": true,
		"ticker": "DNM",
		"price_id": "price id",
		"fee_token": true,
		"gas_price_levels": {
			"low": 0.034,
			"average": 0.05,
			"high": 0.06
		},
		"fetch_price": true,
		"relayer_denom": true,
		"minimum_thresh_relayer_balance": 24000
	}
	],
	"demeris_addresses": ["0x12324", "0x34567"],
	"genesis_hash": "0x123456",
	"node_info": {
		"endpoint": "https://foo.bar:1234",
		"chain_id": "my_chain",
		"bech32_config": {
			"main_prefix": "prefix",
			"prefix_account": "account",
			"prefix_validator": "validator",
			"prefix_consensus": "consensus",
			"prefix_public": "public",
			"prefix_operator": "operator"
		}
	},
	"valid_block_thresh": "32m",
	"derivation_path": "m/44'/0'/0'",
	"supported_wallets": ["Keplr", "Some other"],
	"block_explorer": "https://explorer.com"
}
`
	ChainWithPublicNodeEndpoints = `
{
	"enabled": true,
	"chain_name": "foo",
	"logo": "logo.png",
	"display_name": "FooBar",
	"primary_channel": {
		"key1": "value1",
		"key2": "value2"
	},
	"denoms": [
	{
		"name": "denom1",
		"display_name": "Denom 1",
		"logo": "https://logo.com",
		"precision": 12,
		"verified": true,
		"stakable": true,
		"ticker": "DNM",
		"price_id": "price id",
		"fee_token": true,
		"gas_price_levels": {
			"low": 0.034,
			"average": 0.05,
			"high": 0.06
		},
		"fetch_price": true,
		"relayer_denom": true,
		"minimum_thresh_relayer_balance": 24000
	}
	],
	"demeris_addresses": ["0x12324", "0x34567"],
	"genesis_hash": "0x123456",
	"node_info": {
		"endpoint": "https://foo.bar:1234",
		"chain_id": "my_chain",
		"bech32_config": {
			"main_prefix": "prefix",
			"prefix_account": "account",
			"prefix_validator": "validator",
			"prefix_consensus": "consensus",
			"prefix_public": "public",
			"prefix_operator": "operator"
		}
	},
	"valid_block_thresh": "32m",
	"derivation_path": "m/44'/0'/0'",
	"supported_wallets": ["Keplr", "Some other"],
	"block_explorer": "https://explorer.com",
	"public_node_endpoints": {
		"tendermint_rpc": "https://localhost:1234",
		"cosmos_api": "https://127.0.01:2345"
	}
}
`
)

func requestWithBody(body string) (req *http.Request) {
	req, _ = http.NewRequest("POST", "/", bytes.NewBufferString(body))
	return
}
