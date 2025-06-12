package transaction

import (
	"metaLand/data/model"
	"time"
)

const (
	Pending        = 0
	Success        = 1
	Failure        = 2
	ConfirmFailure = 3

	BountyDepositContractCreated = 1
	BountyDepositAccount         = 2
	CrowdfundingContractCreated  = 3
	CrowdfundingModified         = 4
	CrowdfundingCancelled        = 5
	CrowdfundingRemoved          = 6
	CrowdfundingBought           = 7
	CrowdfundingSold             = 8
	ReceiptSuccess               = 1
	ReceiptFailure               = 0
)

type GetTransaction struct {
	TransactionId uint64 `gorm:"column:id"`
	ChainID       uint64 `gorm:"column:chain_id"`
	TxHash        string `gorm:"column:tx_hash"`
	SourceID      uint64 `gorm:"column:source_id"`
	SourceType    int    `gorm:"column:source_type"`
	RetryTimes    int    `gorm:"column:retry_times"`
}

type Transaction struct {
	model.RelationBase
	ChainID    uint64    `gorm:"column:chain_id;unique_index:chain_tx_uindex" json:"chainID"`
	TxHash     string    `gorm:"column:tx_hash;unique_index:chain_tx_uindex" json:"txHash"`
	TimeStamp  time.Time `gorm:"column:timestamp"`
	Status     int       `gorm:"column:status" json:"status,omitempty"` // 0:Pending 1:Success 2:Failure
	SourceType int       `gorm:"column:source_type" json:"sourceType"`
	SourceID   int64     `gorm:"column:source_id" json:"sourceID"`
	RetryTimes int       `gorm:"column:retry_times" json:"retryTimes"`
}

// TableName the Transaction table name for gorm
func (Transaction) TableName() string {
	return "transaction"
}
