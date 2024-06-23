package token

import (
	"gorm.io/gorm"
)

type TokenRepository interface {
	Save(token *Token) error
	FindByToken(token string) (*Token, error)
	DeleteByUserID(userID string) error
}

type TokenRepositoryImpl struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepositoryImpl {
	return &TokenRepositoryImpl{db: db}
}

func (r *TokenRepositoryImpl) Save(token *Token) error {
	return r.db.Save(token).Error
}

func (r *TokenRepositoryImpl) FindByToken(token string) (*Token, error) {
	var verificationToken Token
	err := r.db.Where("token = ?", token).First(&verificationToken).Error

	return &verificationToken, err
}

func (r *TokenRepositoryImpl) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&Token{}).Error
}
