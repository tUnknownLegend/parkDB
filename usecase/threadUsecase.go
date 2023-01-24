package usecase

import (
	conf "parkDB/config"
	models "parkDB/models"
	"parkDB/repository"

	"strconv"
)

type ThreadUsecaseInterface interface {
	CreatePosts(slugOrID string, posts *models.Posts) (err error)
	GetPost(slugOrID string) (thread *models.Thread, err error)
	UpdatePost(slugOrID string, thread *models.Thread) (err error)
	GetPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error)
	Vote(slugOrID string, vote *models.Vote) (thread *models.Thread, err error)
}

type ThreadUsecase struct {
	threadRepository repository.ThreadRepositoryInterface
	postRepository   repository.PostRepositoryInterface
	userRepository   repository.UserRepositoryInterface
}

func NewThreadUsecase(
	threadRepository repository.ThreadRepositoryInterface,
	postRepository repository.PostRepositoryInterface,
	userRepository repository.UserRepositoryInterface,
) ThreadUsecaseInterface {
	return &ThreadUsecase{threadRepository: threadRepository, postRepository: postRepository, userRepository: userRepository}
}

func (threadUsecase *ThreadUsecase) CreatePosts(slugOrID string, posts *models.Posts) (err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	if errConv != nil {
		thread, err = threadUsecase.threadRepository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.threadRepository.GetThreadByID(int64(id))
	}

	if err != nil {
		err = conf.NotFoundError
		return
	}

	if len(*posts) == 0 {
		return
	}

	if (*posts)[0].Parent != 0 {
		var parentPost *models.Post
		parentPost, _ = threadUsecase.postRepository.GetPostByID((*posts)[0].Parent)
		if parentPost.Thread != thread.ID {
			err = conf.ConflictError
			return
		}
	}
	_, err = threadUsecase.userRepository.GetByNickname((*posts)[0].Author)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	err = threadUsecase.threadRepository.CreatePosts(thread, posts)
	return
}

func (threadUsecase *ThreadUsecase) GetPost(slugOrID string) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUsecase.threadRepository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.threadRepository.GetThreadByID(int64(id))
	}
	if err != nil {
		err = conf.NotFoundError
		return
	}
	return
}

func (threadUsecase *ThreadUsecase) UpdatePost(slugOrID string, thread *models.Thread) (err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var prevThread *models.Thread
	if errConv != nil {
		prevThread, err = threadUsecase.threadRepository.GetThreadBySlug(slugOrID)
	} else {
		prevThread, err = threadUsecase.threadRepository.GetThreadByID(int64(id))
	}

	if err != nil {
		err = conf.NotFoundError
		return
	}

	if thread.Title != "" {
		prevThread.Title = thread.Title
	}
	if thread.Message != "" {
		prevThread.Message = thread.Message
	}

	err = threadUsecase.threadRepository.UpdatePost(prevThread)
	if err != nil {
		return
	}

	*thread = *prevThread

	return
}

func (threadUsecase *ThreadUsecase) GetPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	if errConv != nil {
		thread, err = threadUsecase.threadRepository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.threadRepository.GetThreadByID(int64(id))
	}

	if err != nil {
		err = conf.NotFoundError
		return
	}

	postsSlice := new([]models.Post)
	switch sort {
	case "tree":
		postsSlice, err = threadUsecase.threadRepository.GetPostsTree(thread.ID, limit, since, desc)
	case "parent_tree":
		postsSlice, err = threadUsecase.threadRepository.GetPostsParentTree(thread.ID, limit, since, desc)
	default:
		postsSlice, err = threadUsecase.threadRepository.GetPostsFlat(thread.ID, limit, since, desc)
	}
	if err != nil {
		return
	}
	posts = new(models.Posts)
	if len(*postsSlice) == 0 {
		*posts = []models.Post{}
	} else {
		*posts = *postsSlice
	}

	return
}

func (threadUsecase *ThreadUsecase) Vote(slugOrID string, vote *models.Vote) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)

	if errConv != nil {
		thread, err = threadUsecase.threadRepository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.threadRepository.GetThreadByID(int64(id))
	}

	if err != nil {
		err = conf.NotFoundError
		return
	}

	err = threadUsecase.threadRepository.Vote(thread.ID, vote)
	if err != nil {
		err = conf.NotFoundError
		return
	}
	thread.Votes, err = threadUsecase.threadRepository.GetVotes(thread.ID)

	return
}
