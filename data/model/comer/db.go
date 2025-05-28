package comer

import (
	"metaLand/data/model/tag"

	"gorm.io/gorm"
)

// ListStartups  list startups
func ListComers(db *gorm.DB, comerID uint64, input *ListComerRequest, startups *[]Comer) (total int64, err error) {
	db = db.Where("is_deleted = false")
	if comerID != 0 {
		db = db.Where("comer_id = ?", comerID)
	}
	if input.Keyword != "" {
		db = db.Where("name like ?", "%"+input.Keyword+"%")
	}
	if input.Mode != 0 {
		db = db.Where("mode = ?", input.Mode)
	}
	if err = db.Table("startup").Count(&total).Error; err != nil {
		return
	}
	if total == 0 {
		return
	}
	err = db.Order("created_at DESC").Limit(input.Limit).Offset(input.Offset).Preload("Wallets").Preload("HashTags", "category = ?", tag.Startup).Preload("Members").Preload("Members.Comer").Preload("Members.ComerProfile").Preload("Follows").Find(startups).Error
	return
}

func FindComer(db *gorm.DB, comerID uint64) (comer *Comer, err error) {
	err = db.Where("id = ?", comerID).First(&comer).Error
	return
}

func FindComerByAddress(db *gorm.DB, address string) (comer *Comer, err error) {
	err = db.Where("address = ? AND is_deleted = 0", address).First(&comer).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func InsertComer(db *gorm.DB, comer *Comer) (err error) {
	err = db.Create(comer).Error
	return
}
