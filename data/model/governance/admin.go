package governance

import "gorm.io/gorm"

func CreateGovernanceAdmins(db *gorm.DB, admins []*GovernanceAdmin) error {
	return db.Create(&admins).Error
}

func DeleteAdminsBySettingId(db *gorm.DB, settingId uint64) error {
	return db.Where("setting_id = ?", settingId).Delete(&GovernanceAdmin{}).Error
}
