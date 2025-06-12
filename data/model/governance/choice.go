package governance

import "gorm.io/gorm"

func CreateChoices(db *gorm.DB, choices []*GovernanceChoice) error {
	return db.Create(&choices).Error
}

func GetProposalChoices(db *gorm.DB, proposalId uint64) (choices []*GovernanceChoice, err error) {
	err = db.Model(&GovernanceChoice{}).Where("proposal_id = ?", proposalId).Order("seq_num asc").Find(&choices).Error
	return
}

func DeleteProposalChoices(db *gorm.DB, proposalId uint64) error {
	return db.Model(&GovernanceChoice{}).Where("proposal_id = ?", proposalId).Update("is_deleted", true).Error
}

func GetChoiceByProposalIdAndChoiceId(db *gorm.DB, proposalId, choiceItemId uint64) (choice GovernanceChoice, err error) {
	err = db.Model(&GovernanceChoice{}).
		Where("proposal_id = ? and id = ?", proposalId, choiceItemId).
		Find(&choice).
		Error
	return
}
