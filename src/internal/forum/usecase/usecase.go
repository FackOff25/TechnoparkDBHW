package usecase

import (
	forumRep "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/repository"
	userRep "github.com/FackOff25/TechnoparkDBHW/src/internal/user/repository"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
)

type UseCaseInterface interface {
	CreateForum(forum *models.Forum) error
	SelectForum(slug string) (*models.Forum, error)
	SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error)
}

type useCase struct {
	forumRepository forumRep.RepositoryInterface
	userRepository  userRep.RepositoryInterface
}

func New(forumRepository forumRep.RepositoryInterface, userRepository userRep.RepositoryInterface) UseCaseInterface {
	return &useCase{
		forumRepository: forumRepository,
		userRepository:  userRepository,
	}
}

func (uc *useCase) CreateForum(forum *models.Forum) error {
	selectedUser, err := uc.userRepository.SelectUserByNickName(forum.User)
	if err != nil {
		return err
	}

	existForum, err := uc.forumRepository.SelectForumBySlug(forum.Slug)
	if err != models.ErrNotFound && err != nil {
		return err
	} else if err == nil {
		forum.User = existForum.User
		forum.Posts = existForum.Posts
		forum.Slug = existForum.Slug
		forum.Threads = existForum.Threads
		forum.Title = existForum.Title
		return models.ErrConflict
	}

	forum.User = selectedUser.NickName

	err = uc.forumRepository.CreateForum(forum)
	if err != nil {
		return err
	}

	return nil
}

func (uc *useCase) SelectForum(slug string) (*models.Forum, error) {
	forum, err := uc.forumRepository.SelectForumBySlug(slug)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (uc *useCase) SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error) {
	_, err := uc.forumRepository.SelectForumBySlug(slug)
	if err != nil {
		return nil, err
	}

	users, err := uc.forumRepository.SelectForumUsers(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return users, nil
}
