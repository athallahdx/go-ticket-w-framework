package repository

import (
	"errors"
	"go-ticket/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll(filter domain.UserFilter, page, limit int) ([]*domain.User, int, error) {
	var users []*domain.User
	var total int64

	query := r.db.Model(&domain.User{})

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := filter.SortBy + " " + filter.Order

	err := query.Order(order).Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (r *userRepository) GetAllWithDeleted(filter domain.UserFilter, page, limit int) ([]*domain.User, int, error) {
	var users []*domain.User
	var total int64

	offset := (page - 1) * limit
	query := r.db.Unscoped().Model(&domain.User{})

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortClause := "created_at DESC"
	allowedSortColumns := map[string]bool{"name": true, "email": true, "role": true, "created_at": true, "deleted_at": true}
	if allowedSortColumns[filter.SortBy] {
		order := "asc"
		if filter.Order == "desc" {
			order = "desc"
		}
		sortClause = filter.SortBy + " " + order
	}

	if err := query.Offset(offset).Limit(limit).Order(sortClause).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (r *userRepository) Update(user *domain.User) error {
	result := r.db.Model(user).Updates(domain.User{
		Name:    user.Name,
		Phone:   user.Phone,
		Profile: user.Profile,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) UpdateRole(id int64, role string) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", id).Update("role", role)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or role unchanged")
	}

	return nil
}

func (r *userRepository) Delete(id int64) error {
	result := r.db.Delete(&domain.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("User not found or already deleted")
	}

	return nil
}
