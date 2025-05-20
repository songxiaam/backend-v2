package startup

import (
	"gorm.io/gorm"
	"metaLand/data/model/tag"
)

// ListStartups  list startups
func ListStartups(db *gorm.DB, comerID uint64, input *ListStartupRequest, startups *[]Startup) (total int64, err error) {
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
