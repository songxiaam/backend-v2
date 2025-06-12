package governance

import "gorm.io/gorm"

func CreateProposalVote(db *gorm.DB, vote *GovernanceVote) error {
	return db.Create(vote).Error
}
func DeleteVoteByProposalIdAndVoterComer(db *gorm.DB, proposalId, voterComerId uint64) error {
	return db.Where("proposal_id = ? and voter_comer_id = ?", proposalId, voterComerId).Delete(&GovernanceVote{}).Error
}

func GetVoteRecordsByProposalId(db *gorm.DB, proposalId uint64, page, size int) (votes []GovernanceVote, total int64, err error) {
	offset := (page - 1) * size
	query := db.Model(&GovernanceVote{}).Where("proposal_id = ?", proposalId)
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Order("created_at DESC").Offset(offset).Limit(size).Find(&votes).Error
	return
}

func GetVoteRecordByProposalIdAndComerId(db *gorm.DB, proposalId, comerId uint64) (vote GovernanceVote, err error) {
	err = db.Model(&GovernanceVote{}).
		Where("proposal_id = ? and voter_comer_id = ?", proposalId, comerId).
		Find(&vote).Error
	return
}
