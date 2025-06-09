package crowdfunding

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

func CreateCrowdfunding(db *gorm.DB, c *Crowdfunding) error {
	return db.Create(c).Error
}

func SelectToBeStartedCrowdfundingListWithin1Min(db *gorm.DB) (list []*Crowdfunding, err error) {
	now := time.Now()
	err = db.Model(&Crowdfunding{}).
		Where("is_deleted = false and status =  ? and (start_time between ? and ? or start_time < ? and end_time > ?)",
			Upcoming,
			now.Add(-30*time.Second),
			now.Add(30*time.Second),
			now,
			now).
		Find(&list).Error
	return
}

func UpdateCrowdfundingStatus(db *gorm.DB, crowdfundingId uint64, status CrowdfundingStatus) error {
	return db.Model(&Crowdfunding{}).Where("is_deleted=0 and id = ?", crowdfundingId).Update("status", status).Error
}

func SelectToBeEndedCrowdfundingList(db *gorm.DB) (list []*Crowdfunding, err error) {
	err = db.Model(&Crowdfunding{}).
		Where("is_deleted = false and status not in  ? and end_time <= ?",
			[]CrowdfundingStatus{Ended, Cancelled, OnChainFailure},
			time.Now().Add(time.Second*30)).
		Find(&list).Error
	return
}

func GetCrowdfundingById(db *gorm.DB, crowdfundingId uint64) (entity Crowdfunding, err error) {
	if err = db.Model(&Crowdfunding{}).Where("is_deleted=0 and id = ?", crowdfundingId).Find(&entity).Error; err != nil {
		return
	}
	return
}

func UpdateCrowdfundingContractAddressAndStatus(db *gorm.DB, fundingID uint64, address string, status CrowdfundingStatus) error {
	return db.Model(&Crowdfunding{}).Where("id = ?", fundingID).Updates(map[string]interface{}{"crowdfunding_contract": address, "status": status}).Error
}

func GetIboRateHistoryById(db *gorm.DB, historyId uint64) (history IboRateHistory, err error) {
	err = db.Model(&IboRateHistory{}).Where("id = ?", historyId).Find(&history).Error
	return
}

func UpdateCrowdfundingSwapStatus(db *gorm.DB, id uint64, status CrowdfundingSwapStatus) error {
	return db.Model(&CrowdfundingSwap{}).Where("id = ?", id).Updates(map[string]interface{}{"status": status}).Error
}

func UpdateCrowdfunding(db *gorm.DB, crowdfundingId uint64, request ModifyRequest) error {
	return db.Model(&Crowdfunding{}).Where("is_deleted=0 and id = ?", crowdfundingId).Updates(Crowdfunding{
		SellInfo:    SellInfo{MaxSellPercent: request.MaxSellPercent},
		BuyInfo:     BuyInfo{BuyPrice: request.BuyPrice, MaxBuyAmount: request.MaxBuyAmount},
		SwapPercent: request.SwapPercent,
		EndTime:     request.EndTime,
	}).Error
}

func GetCrowdfundingSwapById(db *gorm.DB, swapId uint64) (swap CrowdfundingSwap, err error) {
	err = db.Model(&CrowdfundingSwap{}).Where("id = ?", swapId).Find(&swap).Error
	return
}

func UpdateCrowdfundingRaiseBalance(tx *gorm.DB, swap CrowdfundingSwap) error {
	var p clause.Expr
	mp := map[string]interface{}{}
	if swap.Access == Invest {
		p = gorm.Expr("raise_balance + ?", swap.BuyTokenAmount)
		// funding, err := GetCrowdfundingById(tx, swap.CrowdfundingId)
		//if err != nil {
		//	return err
		//}
		//if funding.RaiseBalance.Add(swap.BuyTokenAmount).GreaterThanOrEqual(funding.RaiseGoal) {
		//	mp["status"] = Ended
		//}
		mp["raise_balance"] = p
		return tx.Model(&Crowdfunding{}).Where("is_deleted = false and id = ?", swap.CrowdfundingID).Updates(mp).Error
	}
	p = gorm.Expr("raise_balance - ?", swap.BuyTokenAmount)
	mp["raise_balance"] = p
	return tx.Model(&Crowdfunding{}).Where("is_deleted = false and id = ?", swap.CrowdfundingID).Updates(mp).Error
}

func GetCrowdfundingList(db *gorm.DB, page, size int, keyword string) (list []Crowdfunding, total int64, err error) {
	//var list []Crowdfunding
	tx := db.Where("crowdfunding.is_deleted = 0")

	if strings.TrimSpace(keyword) != "" {
		tx = tx.Joins("inner join startup on startup.id = crowdfunding.startup_id and startup.name like ? and startup.is_deleted = false", "%"+strings.TrimSpace(keyword)+"%")
	}
	//if pagination.Mode == 0 {
	//	tx = tx.Where("crowdfunding.status in ?", []CrowdfundingStatus{Upcoming, Live, Ended, Cancelled})
	//} else {
	//	tx = tx.Where("crowdfunding.status = ?", pagination.Mode)
	//}
	if err := tx.Table("crowdfunding").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := tx.Table("crowdfunding").Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return
}
