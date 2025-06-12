package bounty

import "gorm.io/gorm"

func ListBounties(db *gorm.DB, request *ListBountiesRequest, bounties *[]Bounty) (total int64, err error) {
	db = db.Where("is_deleted = false")
	if request.Keyword != "" {
		db = db.Where("title like ?", "%"+request.Keyword+"%")
	}
	if err = db.Table("bounty").Count(&total).Error; err != nil {
		return
	}
	if total == 0 {
		return
	}
	err = db.Debug().Order("created_at DESC").Limit(request.Limit).Offset(request.Offset).Preload("BountyPaymentPeriod", "is_deleted = ?", "0").Preload("Startup", "is_deleted = ?", "0").Find(bounties).Error
	return
}

// GetBounty /**
func GetBounty(db *gorm.DB, id uint64, bounty *Bounty) (err error) {
	//.Preload("BountyApplicant", "is_deleted = ?", "0").Preload("BountyContact", "is_deleted = ?", "0").Preload("BountyDeposit", "is_deleted = ?", "0")
	return db.Debug().Table("bounty").Where("id = ?", id).
		Preload("BountyApplicants", "is_deleted = ?", "0").
		Preload("BountyContacts", "is_deleted = ?", "0").
		Preload("BountyDeposits", "is_deleted = ?", "0").
		Preload("BountyPaymentTerms", "is_deleted = ?", "0").
		First(bounty).Error
}

func UpdateBountyDepositContract(db *gorm.DB, bountyID uint64, depositContract string) error {
	return db.Model(&Bounty{}).Where("id = ?", bountyID).Update("deposit_contract", depositContract).Error
}

func UpdateBountyDepositStatus(db *gorm.DB, bountyID uint64, status int) error {
	return db.Model(&BountyDeposit{}).Where("bounty_id = ?", bountyID).Update("status", status).Error
}
