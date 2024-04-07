package token

import (
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(token *Token) error {
	return r.db.Save(token).Error
}

func (r *TokenRepository) FindByToken(token string) (*Token, error) {
	var verificationToken Token
	err := r.db.Where("token = ?", token).First(&verificationToken).Error
	
	return &verificationToken, err
}

func (r *TokenRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&Token{}).Error
}