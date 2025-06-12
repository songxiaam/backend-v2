package crowdfunding

import (
	"github.com/shopspring/decimal"
	"time"
)

type TransactionHashRequest struct {
	TxHash string `json:"txHash" binding:"required"`
}

type ModifyRequest struct {
	TransactionHashRequest
	SwapPercent  decimal.Decimal `json:"swapPercent,omitempty"`
	BuyPrice     decimal.Decimal `json:"buyPrice,omitempty"`
	MaxBuyAmount decimal.Decimal `json:"maxBuyAmount,omitempty"`
	//BuyTokenSymbol  string          `json:"buyTokenSymbol,omitempty"`
	//SellTokenSymbol string          `json:"sellTokenSymbol,omitempty"`
	MaxSellPercent decimal.Decimal `json:"maxSellPercent,omitempty"`
	EndTime        time.Time       `json:"endTime,omitempty"`
}
