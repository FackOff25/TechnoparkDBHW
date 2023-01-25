package usecase

import (
	"strconv"
	"time"

	forumRep "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/repository"
	postRep "github.com/FackOff25/TechnoparkDBHW/src/internal/post/repository"
	threadRep "github.com/FackOff25/TechnoparkDBHW/src/internal/thread/repository"
	userRep "github.com/FackOff25/TechnoparkDBHW/src/internal/user/repository"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/go-openapi/strfmt"
)

type UseCaseInterface interface {
	CreatePosts(posts []*models.Post, slugOrId string) error
	UpdatePost(post *models.Post) error
	SelectPost(id uint64, related []string) (*models.PostDetails, error)
	SelectThreadPosts(slugOrId string, limit int, since int, desc bool, sort string) ([]*models.Post, error)
}

type useCase struct {
	postRepository   postRep.RepositoryInterface
	userRepository   userRep.RepositoryInterface
	threadRepository threadRep.RepositoryInterface
	forumRepository  forumRep.RepositoryInterface
}

func New(postRepository postRep.RepositoryInterface, userRepository userRep.RepositoryInterface,
	threadRepository threadRep.RepositoryInterface, forumRepository forumRep.RepositoryInterface) UseCaseInterface {
	return &useCase{
		postRepository:   postRepository,
		userRepository:   userRepository,
		threadRepository: threadRepository,
		forumRepository:  forumRepository,
	}
}

func (uc *useCase) CreatePosts(posts []*models.Post, slugOrId string) error {
	var thread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		thread, err = uc.threadRepository.SelectThreadById(id)
		if err != nil {
			return err
		}
	} else {
		thread, err = uc.threadRepository.SelectThreadBySlug(slugOrId)
		if err != nil {
			return err
		}
	}

	if len(posts) == 0 {
		return nil
	}

	for idx := range posts {
		posts[idx].Thread = thread.Id
		posts[idx].Forum = thread.Forum
		_, err = uc.userRepository.SelectUserByNickName(posts[idx].Author)
		if err != nil {
			return err
		}
		if posts[idx].Parent != 0 {
			selectedPost, err := uc.postRepository.SelectPostById(posts[idx].Parent)
			if err == models.ErrNotFound {
				return models.ErrConflict
			} else if err != nil {
				return err
			} else if selectedPost.Thread != posts[idx].Thread {
				return models.ErrConflict
			}
		}
	}

	timeNow := time.Now().UTC()
	for _, post := range posts {
		post.Created = strfmt.DateTime(timeNow)
	}

	err = uc.postRepository.CreatePosts(posts)
	if err != nil {
		return err
	}

	for idx := range posts {
		err = uc.forumRepository.CreateForumUser(posts[idx].Forum, posts[idx].Author)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *useCase) UpdatePost(post *models.Post) error {
	selectedPost, err := uc.postRepository.SelectPostById(post.Id)
	if err != nil {
		return err
	}

	if post.Message == "" || post.Message == selectedPost.Message {
		post.Id = selectedPost.Id
		post.Message = selectedPost.Message
		post.IsEdited = selectedPost.IsEdited
		post.Author = selectedPost.Author
		post.Created = selectedPost.Created
		post.Forum = selectedPost.Forum
		post.Parent = selectedPost.Parent
		post.Thread = selectedPost.Thread
		return nil
	}

	err = uc.postRepository.UpdatePost(post)
	if err != nil {
		return err
	}

	post.IsEdited = true

	return nil
}

func (uc *useCase) SelectPost(id uint64, related []string) (*models.PostDetails, error) {
	postDetails := models.PostDetails{}

	post, err := uc.postRepository.SelectPostById(id)
	if err != nil {
		return nil, err
	}

	postDetails.Post = post

	for _, elem := range related {
		switch elem {
		case "user":
			user, err := uc.userRepository.SelectUserByNickName(post.Author)
			if err != nil {
				return nil, err
			}
			postDetails.User = user
		case "thread":
			thread, err := uc.threadRepository.SelectThreadById(post.Thread)
			if err != nil {
				return nil, err
			}
			postDetails.Thread = thread
		case "forum":
			forum, err := uc.forumRepository.SelectForumBySlug(post.Forum)
			if err != nil {
				return nil, err
			}
			postDetails.Forum = forum
		}
	}

	return &postDetails, nil
}

func (uc *useCase) SelectThreadPosts(slugOrId string, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
	var selectedThread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		selectedThread, err = uc.threadRepository.SelectThreadById(id)
		if err != nil {
			return nil, err
		}
	} else {
		selectedThread, err = uc.threadRepository.SelectThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
	}

	posts, err := uc.postRepository.SelectThreadPosts(selectedThread.Id, limit, since, desc, sort)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
