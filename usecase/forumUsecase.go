package usecase

import (
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/repository"
)

type ForumUsecaseInterface interface {
	CreateForum(forum *models.Forum) (err error)
	GetForumBySlug(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (err error)
	GetUsersOfForum(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetThreadsOfForum(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
}

type ForumUsecase struct {
	forumRepository  repository.ForumRepositoryInterface
	threadRepository repository.ThreadRepositoryInterface
	userRepository   repository.UserRepositoryInterface
}

func NewForumUsecase(forumRepository repository.ForumRepositoryInterface, threadRepository repository.ThreadRepositoryInterface, userRepository repository.UserRepositoryInterface) ForumUsecaseInterface {
	return &ForumUsecase{forumRepository: forumRepository, threadRepository: threadRepository, userRepository: userRepository}
}

func (forumUsecase *ForumUsecase) CreateForum(forum *models.Forum) (err error) {
	user, err := forumUsecase.userRepository.GetByNickname(forum.User)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	oldForum, _ := forumUsecase.forumRepository.GetForumBySlug(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		err = conf.ConflictError
		return
	}

	forum.User = user.Nickname
	err = forumUsecase.forumRepository.CreateForum(forum)
	return
}

func (forumUsecase *ForumUsecase) GetForumBySlug(slug string) (forum *models.Forum, err error) {
	forum, err = forumUsecase.forumRepository.GetForumBySlug(slug)
	if err != nil {
		err = conf.NotFoundError
	}
	return
}

func (forumUsecase *ForumUsecase) CreateThread(thread *models.Thread) (err error) {
	forum, err := forumUsecase.forumRepository.GetForumBySlug(thread.Forum)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	_, err = forumUsecase.userRepository.GetByNickname(thread.Author)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	oldThread, _ := forumUsecase.threadRepository.GetThreadBySlug(thread.Slug)
	if oldThread.Slug != "" {
		*thread = *oldThread
		err = conf.ConflictError
		return
	}

	thread.Forum = forum.Slug
	err = forumUsecase.threadRepository.CreatePost(thread)
	return
}

func (forumUsecase *ForumUsecase) GetUsersOfForum(slug string, limit int, since string, desc bool) (users *models.Users, err error) {
	_, err = forumUsecase.forumRepository.GetForumBySlug(slug)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	usersSlice, err := forumUsecase.forumRepository.GetUsersOfForum(slug, limit, since, desc)
	if err != nil {
		return
	}
	users = new(models.Users)
	if len(*usersSlice) == 0 {
		*users = []models.User{}
	} else {
		*users = *usersSlice
	}

	return
}

func (forumUsecase *ForumUsecase) GetThreadsOfForum(slug string, limit int, since string, desc bool) (threads *models.Threads, err error) {
	forum, err := forumUsecase.forumRepository.GetForumBySlug(slug)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	threadsSlice, err := forumUsecase.forumRepository.GetThreadsofForum(forum.Slug, limit, since, desc)
	if err != nil {
		return
	}
	threads = new(models.Threads)
	if len(*threadsSlice) == 0 {
		*threads = []models.Thread{}
	} else {
		*threads = *threadsSlice
	}

	return
}
