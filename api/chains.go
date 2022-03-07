package api

import (
	"github.com/emerishq/demeris-backend-models/cns"
	"github.com/emerishq/demeris-backend-models/tracelistener"
)

type ChainsResponse struct {
	Chains []SupportedChain `json:"chains"`
}

type SupportedChain struct {
	ChainName   string `json:"chain_name"`
	DisplayName string `json:"display_name"`
	Logo        string `json:"logo"`
}

type ChainResponse struct {
	Chain cns.Chain `json:"chain"`
}

type Bech32ConfigResponse struct {
	Bech32Config cns.Bech32Config `json:"bech32_config"`
}

type FeeResponse struct {
	Denoms cns.DenomList `json:"denoms"`
}

type FeeAddressResponse struct {
	FeeAddress []string `json:"fee_address"`
}

type FeeAddress struct {
	ChainName  string   `json:"chain_name"`
	FeeAddress []string `json:"fee_address"`
}

type FeeAddressesResponse struct {
	FeeAddresses []FeeAddress `json:"fee_addresses"`
}

type FeeTokenResponse struct {
	FeeTokens []cns.Denom `json:"fee_tokens"`
}

type PrimaryChannel struct {
	Counterparty string `json:"counterparty"`
	ChannelName  string `json:"channel_name"`
}

type PrimaryChannelResponse struct {
	Channel PrimaryChannel `json:"primary_channel"`
}

type PrimaryChannelsResponse struct {
	Channels []PrimaryChannel `json:"primary_channels"`
}

type Trace struct {
	Channel          string `json:"channel,omitempty"`
	Port             string `json:"port,omitempty"`
	ClientId         string `json:"client_id,omitempty"`
	ChainName        string `json:"chain_name,omitempty"`
	CounterpartyName string `json:"counterparty_name,omitempty"`
}

type VerifiedTrace struct {
	IbcDenom  string  `json:"ibc_denom,omitempty"`
	BaseDenom string  `json:"base_denom,omitempty"`
	Verified  bool    `json:"verified"`
	Path      string  `json:"path,omitempty"`
	Trace     []Trace `json:"trace,omitempty"`
	Cause     string  `json:"cause,omitempty"`
}

type VerifiedTraceResponse struct {
	VerifiedTrace VerifiedTrace `json:"verify_trace"`
}

type StatusResponse struct {
	Online bool `json:"online"`
}

type NumbersResponse struct {
	Numbers tracelistener.AuthRow `json:"numbers"`
}

type ValidatorsResponse struct {
	Validators []tracelistener.ValidatorRow `json:"validators"`
}

//nolint //deadcode but used in swagger generation
type InflationResponse struct {
	Inflation string `json:"inflation"`
}

// nolint :ditto
type ParamsResponse struct {
	Params struct {
		MintDenom           string `json:"mint_denom"`
		InflationRateChange string `json:"inflation_rate_change"`
		InflationMax        string `json:"inflation_max"`
		InflationMin        string `json:"inflation_min"`
		GoalBonded          string `json:"goal_bonded"`
		BlocksPerYear       string `json:"blocks_per_year"`
	} `json:"params"`
}

// nolint :ditto
type AnnualProvisionsResponse struct {
	AnnualProvisions string `json:"annual_provisions"`
}

type Coin struct {
	Denom  string `json:"denom,omitempty"`
	Amount string `json:"amount,omitempty"`
}

type SupplyResponse struct {
	Supply []Coin `json:"supply,omitempty"`
}
