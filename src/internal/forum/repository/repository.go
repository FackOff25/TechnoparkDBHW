package repository

import (
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
	"strings"
)

type RepositoryInterface interface {
	CreateForum(forum *models.Forum) error
	SelectForumBySlug(slug string) (*models.Forum, error)
	SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error)
	CreateForumUser(forum string, user string) error
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryInterface {
	return &dataBase{
		db: db,
	}
}

func (dbForum *dataBase) CreateForum(forum *models.Forum) error {
	req := dbForum.db.Omit("posts", "threads")

	req.Create(forum)
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table forums")
	}

	return nil
}

func (dbForum *dataBase) SelectForumBySlug(slug string) (*models.Forum, error) {
	forum := models.Forum{}

	req := dbForum.db.Where("slug = ?", slug)

	req.Take(&forum)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table forums")
	}

	return &forum, nil
}

func (dbForum *dataBase) SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error) {
	users := make([]*models.User, 0, 10)

	subReq := dbForum.db.Select("\"user\"").Table("users_forum").Where("forum = ?", slug)
	req := dbForum.db.Where("nickname IN (?)", subReq)

	req.Find(&users)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table forum_user")
	}

	var sortFunc func(i, j int) bool
	if desc {
		sortFunc = func(i, j int) bool {
			return models.UserComp(users[i], users[j])
		}
	} else {
		sortFunc = func(i, j int) bool {
			return models.UserComp(users[j], users[i])
		}
	}
	sort.Slice(users, sortFunc)

	id := 0
	if since != "" {
		for idx, user := range users {
			if strings.ToLower(user.NickName) == strings.ToLower(since) {
				id = idx + 1
				break
			}
		}
	}

	if len(users)-id < limit {
		limit = len(users) - id
	}
	users = users[id : id+limit]

	return users, nil
}

func (dbForum *dataBase) CreateForumUser(forum string, user string) error {
	fu := models.ForumUser{
		Forum: forum,
		User:  user,
	}
	req := dbForum.db.Table("users_forum").Clauses(clause.OnConflict{DoNothing: true})

	req.Create(&fu)
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table forum_user")
	}

	return nil
}
