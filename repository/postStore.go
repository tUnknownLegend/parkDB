package repository

import (
	models "parkDB/models"
	"time"

	"github.com/jackc/pgx"
)

type PostRepositoryInterface interface {
	GetPostByID(id int64) (post *models.Post, err error)
	UpdatePost(post *models.Post) (err error)
}

type PostStore struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) PostRepositoryInterface {
	return &PostStore{db: db}
}

func (postStore *PostStore) GetPostByID(id int64) (post *models.Post, err error) {
	post = &models.Post{}
	postTime := time.Time{}
	err = postStore.db.QueryRow("SELECT id, COALESCE(parent, 0), author, message, edited, forum, thread, created FROM posts "+
		"WHERE id = $1", id).
		Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &postTime)
	post.Created = postTime.Format(time.RFC3339)
	return
}

func (postStore *PostStore) UpdatePost(post *models.Post) (err error) {
	_, err = postStore.db.Exec("UPDATE posts SET message = $1, edited = $2 WHERE id = $3;", post.Message, post.IsEdited, post.ID)
	return
}
