package crowdfunding

import "gorm.io/gorm"

func CreateIboRateHistory(db *gorm.DB, history *IboRateHistory) error {
	return db.Model(&IboRateHistory{}).Create(&history).Error
}
