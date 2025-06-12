package governance

import "gorm.io/gorm"

func GetGovernanceSetting(db *gorm.DB, startupId uint64) (setting GovernanceSetting, err error) {
	err = db.Model(&GovernanceSetting{}).Where("startup_id = ?", startupId).Find(&setting).Error
	return
}

func CreateGovernanceSetting(db *gorm.DB, setting *GovernanceSetting) error {
	return db.Create(setting).Error
}

func UpdateGovernanceSetting(db *gorm.DB, id uint64, updating *GovernanceSetting) error {
	return db.Model(&GovernanceSetting{}).Where("id = ?", id).Updates(updating).Error
}
