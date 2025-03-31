package mysql

import (
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/pkg/db"
)

type adminRepository struct{}

func NewAdminRepository() repository.AdminRepository {
	return &adminRepository{}
}

func (r *adminRepository) GetByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	err := db.DB.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) GetByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := db.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) GetAll() ([]*model.Admin, error) {
	var admins []*model.Admin
	err := db.DB.Find(&admins).Error
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *adminRepository) Create(admin *model.Admin) error {
	return db.DB.Create(admin).Error
}

func (r *adminRepository) Update(admin *model.Admin) error {
	return db.DB.Save(admin).Error
}

func (r *adminRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Admin{}, id).Error
}

func (r *adminRepository) CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := db.DB.Model(&model.Admin{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
