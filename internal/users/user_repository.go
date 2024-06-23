package users

import "gorm.io/gorm"

type UserRepository interface {
	Save(user *User) error
	Delete(id string) error
	FindById(id string) (*User, error)
	FindByEmail(email string) (*User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (rep *UserRepositoryImpl) Save(user *User) error {
	err := rep.db.Save(user).Error

	return err
}

func (rep *UserRepositoryImpl) Delete(id string) error {
	err := rep.db.Delete(&User{}, id).Error

	return err
}

func (rep *UserRepositoryImpl) FindById(id string) (*User, error) {
	var user User

	if err := rep.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (rep *UserRepositoryImpl) FindByEmail(email string) (*User, error) {
	var user User

	if err := rep.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
