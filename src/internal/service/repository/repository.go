package repository

import (
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	ClearData() error
	SelectStatus() (*models.ServiceStatus, error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryInterface {
	return &dataBase{
		db: db,
	}
}

func (dbService *dataBase) ClearData() error {
	req := dbService.db.Exec("TRUNCATE posts, threads, forums, users, users_forum cascade;")
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error")
	}

	return nil
}

func (dbService *dataBase) SelectStatus() (*models.ServiceStatus, error) {
	status := models.ServiceStatus{}

	var count int64
	req := dbService.db.Model(&models.User{}).Count(&count)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table users")
	}
	status.UserCount = count
	req = dbService.db.Model(&models.Forum{}).Count(&count)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table forums")
	}
	status.ForumCount = count
	req = dbService.db.Model(&models.Thread{}).Count(&count)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table posts")
	}
	status.ThreadCount = count
	req = dbService.db.Model(&models.Post{}).Count(&count)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table threads")
	}
	status.PostCount = count

	return &status, nil
}
