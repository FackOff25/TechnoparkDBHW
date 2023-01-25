package repository

import (
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type RepositoryInterface interface {
	CreatePosts(posts []*models.Post) error
	UpdatePost(post *models.Post) error
	SelectPostById(id uint64) (*models.Post, error)
	SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryInterface {
	return &dataBase{
		db: db,
	}
}

func (dbPost *dataBase) CreatePosts(posts []*models.Post) error {
	req := dbPost.db.Create(&posts)
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table posts")
	}

	return nil
}

func (dbPost *dataBase) UpdatePost(post *models.Post) error {
	req := dbPost.db.Model(post).Clauses(clause.Returning{}).Updates(models.Post{Message: post.Message, IsEdited: true})
	if req.Error != nil {
		return errors.Wrap(req.Error, "database error: table posts")
	}

	return nil
}

func (dbPost *dataBase) SelectPostById(id uint64) (*models.Post, error) {
	post := models.Post{}

	req := dbPost.db.Where("id = ?", id)

	req.Take(&post)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table posts")
	}

	return &post, nil
}

func (dbPost *dataBase) SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
	switch sort {
	case "flat":
		return dbPost.SelectThreadPostsSortFlat(id, limit, since, desc)
	case "tree":
		return dbPost.SelectThreadPostsSortTree(id, limit, since, desc)
	case "parent_tree":
		return dbPost.SelectThreadPostsSortParentTree(id, limit, since, desc)
	default:
		return dbPost.SelectThreadPostsSortFlat(id, limit, since, desc)
	}
}

func (dbPost *dataBase) SelectThreadPostsSortFlat(id uint64, limit int, since int, desc bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 10)

	req := dbPost.db.Limit(limit)

	whereInter := "thread = " + strconv.FormatUint(id, 10)

	if since != 0 {
		if desc {
			whereInter += "AND id < "
		} else {
			whereInter += "AND id > "
		}
		whereInter += strconv.Itoa(since)
	}

	req = req.Where(whereInter)

	if desc {
		req = req.Order("id desc")
	} else {
		req = req.Order("id")
	}

	req.Find(&posts)
	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table posts")
	}

	return posts, nil
}

func (dbPost *dataBase) SelectThreadPostsSortTree(id uint64, limit int, since int, desc bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 10)

	req := dbPost.db.Limit(limit)

	whereInter := "thread = " + strconv.FormatUint(id, 10)

	if since != 0 {
		subReq := dbPost.db.Table("posts").Select("post_tree")
		subReq = subReq.Where("id = ?", since)
		if desc {
			whereInter += "AND post_tree < (?)"
		} else {
			whereInter += "AND post_tree > (?)"
		}
		req = req.Where(whereInter, subReq)
	} else {
		req = req.Where(whereInter)
	}

	if desc {
		req = req.Order("post_tree desc")
	} else {
		req = req.Order("post_tree")
	}

	req.Find(&posts)

	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table posts")
	}

	return posts, nil
}

func (dbPost *dataBase) SelectThreadPostsSortParentTree(id uint64, limit int, since int, desc bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 10)

	req := dbPost.db

	whereInter := "post_tree[1] IN (?)"

	subReqWhereInter := "parent = 0 AND thread = ?"
	subReq := dbPost.db.Limit(limit)
	if since != 0 {
		if desc {
			subReqWhereInter += " AND id < (?)"
		} else {
			subReqWhereInter += " AND id > (?)"
		}

		subSubReq := dbPost.db.Table("posts").Select("post_tree[1]")
		subSubReq = subSubReq.Where("id = ?", since)

		subReq = dbPost.db.Table("posts").Limit(limit).Where(subReqWhereInter, id, subSubReq)
	} else {
		subReq = dbPost.db.Table("posts").Limit(limit).Where(subReqWhereInter, id)
	}

	if desc {
		subReq = subReq.Order("id desc")
	} else {
		subReq = subReq.Order("id")
	}
	subReq.Select("id")
	req = req.Where(whereInter, subReq)

	if desc {
		req = req.Order("post_tree[1] desc, post_tree")
	} else {
		req = req.Order("post_tree")
	}

	req.Find(&posts)

	if req.Error != nil {
		return nil, errors.Wrap(req.Error, "database error: table posts")
	}

	return posts, nil
}
