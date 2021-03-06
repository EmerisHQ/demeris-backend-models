package cns

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/mod/semver"
)

// Chain represents CNS chain metadata row on the database.
type Chain struct {
	ID                  uint64              `diff:"-" db:"id" json:"-"`
	Enabled             bool                `diff:"-" db:"enabled" json:"enabled"`                                                                          // boolean that marks whether the given chain is enabled or not (when enabled, API endpoints will return data)
	ChainName           string              `db:"chain_name" binding:"required" json:"chain_name"`                                                          // the unique name of the chain
	Logo                string              `diff:"-" db:"logo" binding:"required" json:"logo"`                                                             // logo of the chain
	DisplayName         string              `diff:"-" db:"display_name" binding:"required" json:"display_name"`                                             // user-friendly chain name
	PrimaryChannel      DbStringMap         `diff:"-" db:"primary_channel"  json:"primary_channel"`                                                         // a mapping of chain name to primary channel
	Denoms              DenomList           `diff:"-" db:"denoms" binding:"dive" json:"denoms"`                                                             // a list of denoms native to the chain
	DemerisAddresses    pq.StringArray      `diff:"-" db:"demeris_addresses" binding:"required" json:"demeris_addresses"`                                   // the addresses on which we accept fee payments
	GenesisHash         string              `diff:"-" db:"genesis_hash" binding:"required" json:"genesis_hash"`                                             // hash of the chain's genesis file
	NodeInfo            NodeInfo            `diff:"-" db:"node_info" binding:"required,dive" json:"node_info"`                                              // info required to query full-node (e.g. to submit tx)
	ValidBlockThresh    Threshold           `diff:"-" db:"valid_block_thresh" binding:"required" json:"valid_block_thresh" swaggertype:"primitive,integer"` // valid block time expressed in time.Duration format
	DerivationPath      string              `diff:"-" db:"derivation_path" binding:"required,derivationpath" json:"derivation_path"`                        // chain derivation path
	SupportedWallets    pq.StringArray      `diff:"-" db:"supported_wallets" binding:"required" json:"supported_wallets"`                                   // the list of supported wallets
	BlockExplorer       string              `diff:"-" db:"block_explorer" json:"block_explorer"`                                                            // block explorer url
	PublicNodeEndpoints PublicNodeEndpoints `diff:"-" db:"public_node_endpoints" binding:"dive" json:"public_node_endpoints,omitempty"`                     // endpoints for non-natively supported chains
	CosmosSDKVersion    string              `diff:"-" db:"cosmos_sdk_version" binding:"required,semver" json:"cosmos_sdk_version,omitempty"`                // Cosmos SDK version used by the chain
}

// VerifiedTokens returns a DenomList of native denoms that are verified.
func (c Chain) VerifiedTokens() DenomList {
	var ret DenomList
	for _, ft := range c.Denoms {
		if !ft.Verified {
			continue
		}

		ret = append(ret, ft)
	}

	return ret
}

// FeeTokens returns a DenomList of denoms that are usable as fee.
func (c Chain) FeeTokens() DenomList {
	var ret DenomList
	for _, ft := range c.Denoms {
		if !ft.FeeToken {
			continue
		}

		ret = append(ret, ft)
	}

	return ret
}

// RelayerToken returns the relayer token for a given chain.
func (c Chain) RelayerToken() Denom {
	for _, ft := range c.Denoms {
		if ft.RelayerDenom {
			return ft
		}
	}

	panic("relayer token not defined")
}

func (c Chain) MajorSDKVersion() string {
	rawVersion := semver.MajorMinor(c.CosmosSDKVersion)
	return strings.Split(rawVersion, ".")[1]
}

// Threshold is a database-friendly time.Duration.
type Threshold time.Duration

func (t *Threshold) UnmarshalJSON(bytes []byte) error {
	str := ""

	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*t = Threshold(d)

	return nil
}

func (t Threshold) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Duration().String())
}

// Duration returns t as time.Duration.
func (t Threshold) Duration() time.Duration {
	return time.Duration(t)
}

// Scan is the sql.Scanner implementation for Threshold.
func (t *Threshold) Scan(value interface{}) error {
	vs, ok := value.(string)
	if !ok {
		return fmt.Errorf("threshold value is of type %T, not string", value)
	}

	vsd, err := time.ParseDuration(vs)
	if err != nil {
		return fmt.Errorf("cannot parse value as duration, %w", err)
	}

	*t = Threshold(vsd)

	return nil
}

// Value is the driver.Value implementation for Threshold.
func (t Threshold) Value() (driver.Value, error) {
	td := time.Duration(t)
	return driver.Value(td.String()), nil
}

// NodeInfo holds information useful to connect to a full node and broadcast transactions.
type NodeInfo struct {
	Endpoint     string       `binding:"required" json:"endpoint"`
	ChainID      string       `binding:"required" json:"chain_id"`
	Bech32Config Bech32Config `binding:"required,dive" json:"bech32_config"`
}

// Scan is the sql.Scanner implementation for DbStringMap.
func (a *NodeInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

// PublicNodeEndpoints holds information for experimental chains, i.e. not natively supported by our wallets.
// This enables the "Suggest Chain" feature in the front-end
type PublicNodeEndpoints struct {
	TendermintRPC []string `binding:"required_with=CosmosAPI,omitempty,min=1,dive,cosmosrpcurl" json:"tendermint_rpc"`
	CosmosAPI     []string `binding:"required_with=TendermintRPC,omitempty,min=1,dive,cosmosrpcurl" json:"cosmos_api"`
}

// Scan is the sql.Scanner implementation for PublicNodeEndpoints.
func (a *PublicNodeEndpoints) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

// GasPrice holds gas prices.
type GasPrice struct {
	Low     float64 `json:"low"`
	Average float64 `json:"average"`
	High    float64 `json:"high"`
}

func (a GasPrice) Empty() bool {
	return a == GasPrice{}
}

// Scan is the sql.Scanner implementation for DbStringMap.
func (a *GasPrice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

// Bech32Config represents the chain's bech32 configuration
type Bech32Config struct {
	MainPrefix      string `json:"main_prefix" binding:"required"`
	PrefixAccount   string `json:"prefix_account" binding:"required"`
	PrefixValidator string `json:"prefix_validator" binding:"required"`
	PrefixConsensus string `json:"prefix_consensus" binding:"required"`
	PrefixPublic    string `json:"prefix_public" binding:"required"`
	PrefixOperator  string `json:"prefix_operator" binding:"required"`
}

// MarshalJSON implements the json.Marshaler interface.
// Returns the json representation of Bech32Config with prefixes methods as fields.
func (b Bech32Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(bech32ConfigMarshaled{
		MainPrefix:      b.MainPrefix,
		PrefixAccount:   b.PrefixAccount,
		PrefixValidator: b.PrefixValidator,
		PrefixConsensus: b.PrefixConsensus,
		PrefixPublic:    b.PrefixPublic,
		PrefixOperator:  b.PrefixOperator,
		AccAddr:         b.Bech32PrefixAccAddr(),
		AccPub:          b.Bech32PrefixAccPub(),
		ValAddr:         b.Bech32PrefixValAddr(),
		ValPub:          b.Bech32PrefixValPub(),
		ConsAddr:        b.Bech32PrefixConsAddr(),
		ConsPub:         b.Bech32PrefixConsPub(),
	})
}

// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
func (b Bech32Config) Bech32PrefixAccAddr() string {
	return b.MainPrefix
}

// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
func (b Bech32Config) Bech32PrefixAccPub() string { return b.MainPrefix + b.PrefixPublic }

// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
func (b Bech32Config) Bech32PrefixValAddr() string {
	return b.MainPrefix + b.PrefixValidator + b.PrefixOperator
}

// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
func (b Bech32Config) Bech32PrefixValPub() string {
	return b.MainPrefix + b.PrefixValidator + b.PrefixOperator + b.PrefixPublic
}

// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
func (b Bech32Config) Bech32PrefixConsAddr() string {
	return b.MainPrefix + b.PrefixValidator + b.PrefixConsensus
}

// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
func (b Bech32Config) Bech32PrefixConsPub() string {
	return b.MainPrefix + b.PrefixValidator + b.PrefixConsensus + b.PrefixPublic
}

type bech32ConfigMarshaled struct {
	MainPrefix      string `json:"main_prefix" binding:"required"`
	PrefixAccount   string `json:"prefix_account" binding:"required"`
	PrefixValidator string `json:"prefix_validator" binding:"required"`
	PrefixConsensus string `json:"prefix_consensus" binding:"required"`
	PrefixPublic    string `json:"prefix_public" binding:"required"`
	PrefixOperator  string `json:"prefix_operator" binding:"required"`
	AccAddr         string `json:"acc_addr,omitempty" db:"-"`
	AccPub          string `json:"acc_pub,omitempty" db:"-"`
	ValAddr         string `json:"val_addr,omitempty" db:"-"`
	ValPub          string `json:"val_pub,omitempty" db:"-"`
	ConsAddr        string `json:"cons_addr,omitempty" db:"-"`
	ConsPub         string `json:"cons_pub,omitempty" db:"-"`
}

// Denom holds a token denomination and its verification status.
type Denom struct {
	Name                        string   `db:"name" binding:"required" json:"name,omitempty"`
	DisplayName                 string   `db:"display_name" json:"display_name"`
	Logo                        string   `db:"logo" json:"logo,omitempty"`
	Precision                   int64    `db:"precision" json:"precision,omitempty"`
	Verified                    bool     `db:"verified" json:"verified,omitempty"`
	Stakable                    bool     `db:"stakable" json:"stakable,omitempty"`
	Ticker                      string   `db:"ticker" json:"ticker,omitempty"`
	PriceID                     string   `db:"price_id" json:"price_id,omitempty"`
	FeeToken                    bool     `db:"fee_token" json:"fee_token,omitempty"`
	GasPriceLevels              GasPrice `db:"gas_price_levels" json:"gas_price_levels"`
	FetchPrice                  bool     `db:"fetch_price" json:"fetch_price"`
	RelayerDenom                bool     `db:"relayer_denom" json:"relayer_denom"`
	MinimumThreshRelayerBalance *int64   `db:"minimum_thresh_relayer_balance" json:"minimum_thresh_relayer_balance,omitempty"`
}

// DenomList represents a slice of Denom.
type DenomList []Denom

// Scan is the sql.Scanner implementation for DenomList.
func (a *DenomList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

// DbStringMap represent a JSON database-enabled string map.
type DbStringMap map[string]string

// Scan is the sql.Scanner implementation for DbStringMap.
func (a *DbStringMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return nil
}

// ChannelQuery represents a query to get a specified channel or counterparty data.
type ChannelQuery struct {
	ChainName    string `db:"chain_name" json:"chain_name"`
	Counterparty string `db:"key" json:"counterparty"`
	ChannelName  string `db:"value" json:"channel_name"`
}
