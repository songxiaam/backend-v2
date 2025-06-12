package governance

import "gorm.io/gorm"

func CreateGovernanceAdmins(db *gorm.DB, admins []*GovernanceAdmin) error {
	return db.Create(&admins).Error
}

func DeleteAdminsBySettingId(db *gorm.DB, settingId uint64) error {
	return db.Where("setting_id = ?", settingId).Delete(&GovernanceAdmin{}).Error
}

func GetGovernanceAdminsByStartupId(db *gorm.DB, startupId uint64) (admins []*GovernanceAdmin, err error) {
	err = db.Table("governance_admin").
		Joins("left join governance_setting on governance_setting.id = governance_admin.setting_id").
		Where("governance_setting.startup_id = ?", startupId).
		Scan(&admins).Error
	return
}

func GetGovernanceAdmins(db *gorm.DB, settingId uint64) (admins []*GovernanceAdmin, err error) {
	err = db.Model(&GovernanceAdmin{}).Where("setting_id = ?", settingId).Find(&admins).Error
	return
}
