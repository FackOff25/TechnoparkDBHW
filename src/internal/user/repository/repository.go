package repository

import (
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryInterface interface {
	SelectUserByNickName(nickname string) (*models.User, error)
	SelectUserByEmail(email string) (*models.User, error)
	SelectUsersByNickNameOrEmail(nickname string, email string) ([]*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryInterface {
	return &dataBase{
		db: db,
	}
}

func (dbUser *dataBase) CreateUser(user *models.User) error {
	req := dbUser.db.Create(user)
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table users")
	}

	return nil
}

func (dbUser *dataBase) SelectUserByNickName(nickname string) (*models.User, error) {
	user := models.User{}

	req := dbUser.db.Where("nickname = ?", nickname)

	req.Take(&user)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table users")
	}

	return &user, nil
}

func (dbUser *dataBase) SelectUsersByNickNameOrEmail(nickname string, email string) ([]*models.User, error) {
	users := make([]*models.User, 0, 10)

	req := dbUser.db.Where("email = ? OR nickname = ?", email, nickname)

	req.Find(&users)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table users")
	}

	return users, nil
}

func (dbUser *dataBase) SelectUserByEmail(email string) (*models.User, error) {
	user := models.User{}

	req := dbUser.db.Where("email = ?", email)

	req.Take(&user)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table users")
	}

	return &user, nil
}

func (dbUser *dataBase) UpdateUser(user *models.User) error {
	req := dbUser.db.Model(user).Clauses(clause.Returning{}).Updates(models.User{About: user.About, Email: user.Email, FullName: user.FullName})
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table users")
	}

	return nil
}
