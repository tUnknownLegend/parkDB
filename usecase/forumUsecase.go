package usecase

import (
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/repository"
)

type ForumUsecaseInterface interface {
	CreateForum(forum *models.Forum) (err error)
	Get(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (err error)
	GetUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
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

	oldForum, err := forumUsecase.forumRepository.GetBySlug(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		err = conf.ConflictError
		return
	}

	forum.User = user.Nickname
	err = forumUsecase.forumRepository.Create(forum)
	return
}

func (forumUsecase *ForumUsecase) Get(slug string) (forum *models.Forum, err error) {
	forum, err = forumUsecase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = conf.NotFoundError
	}
	return
}

func (forumUsecase *ForumUsecase) CreateThread(thread *models.Thread) (err error) {
	forum, err := forumUsecase.forumRepository.GetBySlug(thread.Forum)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	_, err = forumUsecase.userRepository.GetByNickname(thread.Author)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	oldThread, err := forumUsecase.threadRepository.GetBySlug(thread.Slug)
	if oldThread.Slug != "" {
		*thread = *oldThread
		err = conf.ConflictError
		return
	}

	thread.Forum = forum.Slug
	err = forumUsecase.threadRepository.Create(thread)
	return
}

func (forumUsecase *ForumUsecase) GetUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error) {
	_, err = forumUsecase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	usersSlice, err := forumUsecase.forumRepository.GetUsers(slug, limit, since, desc)
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

func (forumUsecase *ForumUsecase) GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error) {
	forum, err := forumUsecase.forumRepository.GetBySlug(slug)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	threadsSlice, err := forumUsecase.forumRepository.GetThreads(forum.Slug, limit, since, desc)
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
