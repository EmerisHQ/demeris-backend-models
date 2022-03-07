package api

import "github.com/emerishq/emeris-utils/exported/sdktypes"

type SwapFeesResponse struct {
	Fees sdktypes.Coins `json:"fees"`
}
