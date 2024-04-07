package user

import "gorm.io/gorm"

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (rep *UserRepository) Save(user *User) error {
	err := rep.db.Save(user).Error

	return err
}

func (rep *UserRepository) Delete(id string) error {
	err := rep.db.Delete(&User{}, id).Error

	return err
}

func (rep *UserRepository) FindById(id string) (*User, error) {
	var user User

	if err := rep.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (rep *UserRepository) FindByEmail(email string) (*User, error) {
	var user User

	if err := rep.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
