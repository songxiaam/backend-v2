package governance

import "gorm.io/gorm"

func CreateGovernanceStrategies(db *gorm.DB, strategies []*GovernanceStrategy) error {
	return db.Create(&strategies).Error
}

func DeleteStrategiesBySettingId(db *gorm.DB, settingId uint64) error {
	return db.Where("setting_id = ?", settingId).Delete(&GovernanceStrategy{}).Error
}
