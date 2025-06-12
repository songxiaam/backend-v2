package transaction

import "gorm.io/gorm"

func GetPendingTransactions(db *gorm.DB) (transactionResponse []*GetTransaction, err error) {
	err = db.Table("transaction").Select("id, chain_id, tx_hash, source_id, source_type, retry_times").
		Where("status = 0").Order("created_at asc").Find(&transactionResponse).Error
	if err != nil {
		return nil, err
	}
	return transactionResponse, nil
}

func UpdateTransactionStatusById(db *gorm.DB, transactionId uint64, status int) error {
	return db.Model(&Transaction{}).Where("id = ?", transactionId).
		Updates(map[string]interface{}{"status": status, "retry_times": gorm.Expr("retry_times + 1")}).Error
}

func UpdateTransactionStatus(db *gorm.DB, sourceID uint64, status int) error {
	return db.Model(&Transaction{}).Where("source_id = ?", sourceID).Update("status", status).Error
}

func CreateTransaction(db *gorm.DB, transaction *Transaction) error {
	return db.Create(&transaction).Error
}

func UpdateTransactionStatusWithRetry(db *gorm.DB, sourceID uint64, sourceType, status int) error {
	return db.Model(&Transaction{}).Where("source_id = ? and source_type = ?", sourceID, sourceType).
		Updates(map[string]interface{}{"status": status, "retry_times": gorm.Expr("retry_times + 1")}).Error
}
