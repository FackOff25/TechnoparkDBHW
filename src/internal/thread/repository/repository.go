package repository

import (
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryInterface interface {
	CreateThread(thread *models.Thread) error
	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error)
	SelectThreadById(id uint64) (*models.Thread, error)
	UpdateThread(thread *models.Thread) error
	CreateVote(vote *models.Vote) error
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryInterface {
	return &dataBase{
		db: db,
	}
}

func (dbThread *dataBase) CreateThread(thread *models.Thread) error {
	req := dbThread.db

	if thread.Slug == "" {
		req = req.Omit("votes", "slug")
	} else {
		req = req.Omit("votes")
	}

	req.Create(thread)
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table threads")
	}

	return nil
}

func (dbThread *dataBase) UpdateThread(thread *models.Thread) error {
	req := dbThread.db.Model(thread).Clauses(clause.Returning{}).Updates(models.Thread{Message: thread.Message, Title: thread.Title})
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table threads")
	}

	return nil
}

func (dbThread *dataBase) SelectThreadBySlug(slug string) (*models.Thread, error) {
	thread := models.Thread{}

	req := dbThread.db.Where("slug = ?", slug)

	req.Take(&thread)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table threads")
	}

	return &thread, nil
}

func (dbThread *dataBase) SelectThreadById(id uint64) (*models.Thread, error) {
	thread := models.Thread{}

	req := dbThread.db.Where("id = ?", id)

	req.Take(&thread)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table threads")
	}

	return &thread, nil
}

func (dbThread *dataBase) SelectForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error) {
	threads := make([]*models.Thread, 0, 10)

	req := dbThread.db.Limit(limit)

	if since != "" {
		if desc {
			req = req.Where("forum = ? AND created <= ?", slug, since)
		} else {
			req = req.Where("forum = ? AND created >= ?", slug, since)
		}
	} else {
		req = req.Where("forum = ?", slug)
	}

	if desc {
		req = req.Order("created desc")
	} else {
		req = req.Order("created")
	}

	req.Find(&threads)

	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table threads")
	}

	return threads, nil
}

func (dbThread *dataBase) CreateVote(vote *models.Vote) error {
	req := dbThread.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "thread_id"}, {Name: "nickname"}},
		DoUpdates: clause.AssignmentColumns([]string{"voice"}),
	})

	req.Create(vote)
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table votes")
	}

	return nil
}
