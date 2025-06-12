package governance

import (
	"gorm.io/gorm"
	"time"
)

func GetProposalList(db *gorm.DB, page int, size int, keyword string) (proposals []GovernanceProposal, total int64, err error) {
	offset := (page - 1) * size
	query := db.Model(&GovernanceProposal{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Order("created_at DESC").Offset(offset).Limit(size).Find(&proposals).Error
	return
}

func CreateProposal(db *gorm.DB, proposal *GovernanceProposal) error {
	return db.Create(proposal).Error
}

func GetProposalById(db *gorm.DB, proposalId uint64) (proposal GovernanceProposal, err error) {
	err = db.Model(&GovernanceProposal{}).Where("id = ? and is_deleted = false", proposalId).Find(&proposal).Error
	return
}

func DeleteProposal(db *gorm.DB, proposalId uint64) error {
	return db.Model(&GovernanceProposal{}).Where("id = ?", proposalId).Update("is_deleted", true).Error
}

func UpdateProposalStatus(db *gorm.DB, proposalId uint64, status ProposalStatus) error {
	return db.Model(&GovernanceProposal{}).Where("id = ?", proposalId).Update("status", status).Error
}
func SelectToBeStartedProposalListWithin1Min(db *gorm.DB) (list []GovernanceProposal, err error) {
	now := time.Now()
	err = db.Model(&GovernanceProposal{}).
		Where("status = ? and (start_time between ? and ? or start_time < ? and end_time > ?)",
			ProposalUpcoming,
			now.Add(-30*time.Second),
			now.Add(30*time.Second),
			now,
			now).
		Find(&list).Error
	return
}

func SelectToEndedProposalListWithin1Min(db *gorm.DB) (list []GovernanceProposal, err error) {
	err = db.Model(&GovernanceProposal{}).
		Where("is_deleted = false and status not in  ? and end_time <= ?",
			[]ProposalStatus{ProposalInvalid, ProposalEnded, ProposalUpcoming, ProposalPending},
			time.Now().Add(time.Second*30)).
		Find(&list).Error
	return
}

func GetProposalPublicInfo(db *gorm.DB, proposalId uint64) (proposal ProposalPublicInfo, err error) {
	err = db.
		Where("governance_proposal.id = ?", proposalId).
		Select("governance_proposal.*,governance_setting.*,startup.logo as startup_logo, startup.name as startup_name, comer_profile.avatar as author_comer_avatar, comer_profile.name as author_comer_name").
		Joins("left join startup on startup.id = governance_proposal.startup_id").
		Joins("left join governance_setting on governance_setting.startup_id = governance_proposal.startup_id").
		Joins("left join comer_profile on comer_profile.comer_id = governance_proposal.author_comer_id").
		Table("governance_proposal").Scan(&proposal).Error
	return
}
