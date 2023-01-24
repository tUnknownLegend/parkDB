package usecase

import (
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/repository"
)

type PostUsecaseInterface interface {
	Get(postID int64, relatedData *[]string) (postFull *models.PostFull, err error)
	Update(post *models.Post) (err error)
}

type PostUsecase struct {
	postRepository   repository.PostRepositoryInterface
	userRepository   repository.UserRepositoryInterface
	threadRepository repository.ThreadRepositoryInterface
	forumRepository  repository.ForumRepositoryInterface
}

func NewPostUsecase(
	postRepository repository.PostRepositoryInterface,
	userRepository repository.UserRepositoryInterface,
	threadRepository repository.ThreadRepositoryInterface,
	forumRepository repository.ForumRepositoryInterface,
) PostUsecaseInterface {
	return &PostUsecase{
		postRepository:   postRepository,
		userRepository:   userRepository,
		threadRepository: threadRepository,
		forumRepository:  forumRepository,
	}
}

func (postUsecase *PostUsecase) Get(postID int64, relatedData *[]string) (postFull *models.PostFull, err error) {
	postFull = new(models.PostFull)
	var post *models.Post
	post, err = postUsecase.postRepository.GetPostByID(postID)
	if err != nil {
		err = conf.NotFoundError
		return
	}
	postFull.Post = post

	for _, data := range *relatedData {
		switch data {
		case "user":
			var author *models.User
			author, err = postUsecase.userRepository.GetByNickname(postFull.Post.Author)
			if err != nil {
				err = conf.NotFoundError
			}
			postFull.Author = author
		case "forum":
			var forum *models.Forum
			forum, err = postUsecase.forumRepository.GetForumBySlug(postFull.Post.Forum)
			if err != nil {
				err = conf.NotFoundError
			}
			postFull.Forum = forum
		case "thread":
			var thread *models.Thread
			thread, err = postUsecase.threadRepository.GetThreadByID(postFull.Post.Thread)
			if err != nil {
				err = conf.NotFoundError
			}
			postFull.Thread = thread
		}
	}
	return
}

func (postUsecase *PostUsecase) Update(post *models.Post) (err error) {
	oldPost, err := postUsecase.postRepository.GetPostByID(post.ID)
	if err != nil {
		err = conf.NotFoundError
		return
	}

	if post.Message != "" {
		if oldPost.Message != post.Message {
			oldPost.IsEdited = true
		}
		oldPost.Message = post.Message

		err = postUsecase.postRepository.UpdatePost(oldPost)
		if err != nil {
			return
		}
	}

	*post = *oldPost

	return
}
