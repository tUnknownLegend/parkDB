package repository

import (
	models "parkDB/models"

	"github.com/jackc/pgx"
)

type ServiceRepositoryInterface interface {
	Clear() (err error)
	GetStatus() (status *models.Status, err error)
}

type ServiceStore struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) ServiceRepositoryInterface {
	return &ServiceStore{db: db}
}

func (serviceStore *ServiceStore) Clear() (err error) {
	_, err = serviceStore.db.Exec("TRUNCATE TABLE forums, posts, threads, user_forum, users, votes CASCADE;")
	return
}

func (serviceStore *ServiceStore) GetStatus() (status *models.Status, err error) {
	status = &models.Status{}
	err = serviceStore.db.QueryRow("SELECT (SELECT count(*) FROM users) AS users, "+
		"(SELECT count(*) FROM forums) AS forums, "+
		"(SELECT count(*) FROM threads) AS threads, "+
		"(SELECT count(*) FROM posts) AS posts;").
		Scan(&status.User, &status.Forum, &status.Thread, &status.Post)
	return
}
