package repository

import (
	models "parkDB/models"

	"github.com/jackc/pgx"
)

type ForumRepositoryInterface interface {
	CreateForum(forum *models.Forum) (err error)
	GetForumBySlug(slug string) (forum *models.Forum, err error)
	GetUsersOfForum(slug string, limit int, since string, desc bool) (users *[]models.User, err error)
	GetThreadsofForum(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error)
}

type ForumStore struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) ForumRepositoryInterface {
	return &ForumStore{db: db}
}

func (forumStore *ForumStore) CreateForum(forum *models.Forum) (err error) {
	_, err = forumStore.db.Exec("INSERT INTO forums (title, forumCreator, slug) VALUES ($1, $2, $3);",
		forum.Title, forum.User, forum.Slug)
	return
}

func (forumStore *ForumStore) GetForumBySlug(slug string) (forum *models.Forum, err error) {
	forum = new(models.Forum)
	err = forumStore.db.QueryRow("SELECT title, forumCreator, slug, posts, threads FROM forums WHERE slug = $1;", slug).
		Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return
}

func (forumStore *ForumStore) GetUsersOfForum(slug string, limit int, since string, desc bool) (users *[]models.User, err error) {
	var usersSlice []models.User

	var resultRows *pgx.Rows

	query := "SELECT users.nickname, users.fullname, users.about, users.email FROM users " +
		"LEFT JOIN userForum ON users.nickname = userForum.nickname WHERE userForum.forum = $1"

	if since != "" {
		if desc {
			query += " AND users.nickname < $2 ORDER BY users.nickname DESC"
		} else {
			query += " AND users.nickname > $2 ORDER BY users.nickname"
		}
		query += " LIMIT $3;"
		resultRows, err = forumStore.db.Query(query, slug, since, limit)
	} else {
		if desc {
			query += " ORDER BY users.nickname DESC"
		} else {
			query += " ORDER BY users.nickname"
		}
		query += " LIMIT $2;"
		resultRows, err = forumStore.db.Query(query, slug, limit)
	}

	if err != nil {
		return
	}
	defer resultRows.Close()

	for resultRows.Next() {
		user := models.User{}
		err = resultRows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return
		}
		usersSlice = append(usersSlice, user)
	}
	return &usersSlice, nil
}

func (forumStore *ForumStore) GetThreadsofForum(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error) {
	var threadsSlice []models.Thread

	var resultRows *pgx.Rows

	query := "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1"

	if since != "" {
		if desc {
			query += " AND created <= $2 ORDER BY created DESC"
		} else {
			query += " AND created >= $2 ORDER BY created ASC"
		}
		query += " LIMIT $3;"
		resultRows, err = forumStore.db.Query(query, slug, since, limit)
	} else {
		if desc {
			query += " ORDER BY created DESC"
		} else {
			query += " ORDER BY created ASC"
		}
		query += " LIMIT $2;"
		resultRows, err = forumStore.db.Query(query, slug, limit)
	}

	if err != nil {
		return
	}
	defer resultRows.Close()

	for resultRows.Next() {
		thread := models.Thread{}
		err = resultRows.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return
		}
		threadsSlice = append(threadsSlice, thread)
	}
	return &threadsSlice, nil
}
