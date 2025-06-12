package crowdfunding

import "gorm.io/gorm"

func QuerySwapListByCrowdfundingId(db *gorm.DB, crowdfundingId uint64) (swaps []CrowdfundingSwap, err error) {
	err = db.Model(&CrowdfundingSwap{}).Where("crowdfunding_id = ? and status = ?", crowdfundingId, SwapSuccess).Find(&swaps).Error
	return
}
