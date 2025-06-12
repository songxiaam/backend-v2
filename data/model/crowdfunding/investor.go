package crowdfunding

import "gorm.io/gorm"

func FirstOrCreateInvestor(tx *gorm.DB, crowdfundingId, comerId uint64) (investor Investor, err error) {
	err = tx.Debug().Where(&Investor{CrowdfundingId: crowdfundingId, ComerId: comerId}).FirstOrCreate(&investor).Error
	return
}

func UpdateCrowdfundingInvestor(db *gorm.DB, investor Investor) error {
	return db.Model(&Investor{}).Where("id = ?", investor.ID).Updates(investor).Error
}

func SelectInvestorByCrowdfundingIdAndComerId(db *gorm.DB, crowdfundingId, comerId uint64) (investor Investor, err error) {
	err = db.Model(&Investor{}).Where("crowdfunding_id = ? and comer_id = ?", crowdfundingId, comerId).Find(&investor).Error
	return
}

func CountByCrowdfundingId(db *gorm.DB, crowdfundingId uint64) (count int64, err error) {
	err = db.Model(&Investor{}).Where("crowdfunding_id = ?", count).Count(&count).Error
	return
}
