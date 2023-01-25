package usecase

import (
	forumRep "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/repository"
	threadRep "github.com/FackOff25/TechnoparkDBHW/src/internal/thread/repository"
	userRep "github.com/FackOff25/TechnoparkDBHW/src/internal/user/repository"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"strconv"
)

type UseCaseInterface interface {
	CreateThread(thread *models.Thread) error
	SelectForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error)
	SelectThread(slugOrId string) (*models.Thread, error)
	UpdateThread(thread *models.Thread, slugOrId string) error
	CreateVote(vote *models.Vote, slugOrId string) (*models.Thread, error)
}

type useCase struct {
	forumRepository  forumRep.RepositoryInterface
	threadRepository threadRep.RepositoryInterface
	userRepository   userRep.RepositoryInterface
}

func New(threadRepository threadRep.RepositoryInterface, userRepository userRep.RepositoryInterface, forumRepository forumRep.RepositoryInterface) UseCaseInterface {
	return &useCase{
		threadRepository: threadRepository,
		userRepository:   userRepository,
		forumRepository:  forumRepository,
	}
}

func (uc *useCase) CreateThread(thread *models.Thread) error {
	_, err := uc.userRepository.SelectUserByNickName(thread.Author)
	if err != nil {
		return err
	}

	selectedForum, err := uc.forumRepository.SelectForumBySlug(thread.Forum)
	if err != nil {
		return err
	}

	if thread.Slug != "" {
		existThread, err := uc.threadRepository.SelectThreadBySlug(thread.Slug)
		if err != models.ErrNotFound && err != nil {
			return err
		} else if err == nil {
			thread.Id = existThread.Id
			thread.Author = existThread.Author
			thread.Created = existThread.Created
			thread.Forum = existThread.Forum
			thread.Message = existThread.Message
			thread.Slug = existThread.Slug
			thread.Title = existThread.Title
			thread.Votes = existThread.Votes
			return models.ErrConflict
		}
	}

	thread.Forum = selectedForum.Slug

	err = uc.threadRepository.CreateThread(thread)
	if err != nil {
		return err
	}

	err = uc.forumRepository.CreateForumUser(thread.Forum, thread.Author)
	if err != nil {
		return err
	}

	return nil
}

func (uc *useCase) SelectForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error) {
	_, err := uc.forumRepository.SelectForumBySlug(slug)
	if err != nil {
		return nil, err
	}

	threads, err := uc.threadRepository.SelectForumThreads(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (uc *useCase) SelectThread(slugOrId string) (*models.Thread, error) {
	var thread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		thread, err = uc.threadRepository.SelectThreadById(id)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = uc.threadRepository.SelectThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
	}

	return thread, nil
}

func (uc *useCase) UpdateThread(thread *models.Thread, slugOrId string) error {
	var selectedThread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		selectedThread, err = uc.threadRepository.SelectThreadById(id)
		if err != nil {
			return err
		}
	} else {
		selectedThread, err = uc.threadRepository.SelectThreadBySlug(slugOrId)
		if err != nil {
			return err
		}
		id = selectedThread.Id
	}

	if thread.Title == "" && thread.Message == "" {
		thread.Author = selectedThread.Author
		thread.Created = selectedThread.Created
		thread.Forum = selectedThread.Forum
		thread.Id = selectedThread.Id
		thread.Slug = selectedThread.Slug
		thread.Votes = selectedThread.Votes
		thread.Title = selectedThread.Title
		thread.Message = selectedThread.Message
		return nil
	}

	thread.Id = id

	err = uc.threadRepository.UpdateThread(thread)
	if err != nil {
		return err
	}

	return nil
}

func (uc *useCase) CreateVote(vote *models.Vote, slugOrId string) (*models.Thread, error) {
	_, err := uc.userRepository.SelectUserByNickName(vote.NickName)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err != nil {
		thread, err := uc.threadRepository.SelectThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
		id = thread.Id
	} else {
		_, err := uc.threadRepository.SelectThreadById(id)
		if err != nil {
			return nil, err
		}
	}

	vote.ThreadId = id

	err = uc.threadRepository.CreateVote(vote)
	if err != nil {
		return nil, err
	}

	thread, err := uc.threadRepository.SelectThreadById(id)
	if err != nil {
		return nil, err
	}

	return thread, nil
}
