package postgre

import (
	"InternProj/internal/models"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreStore struct {
	db *pgxpool.Pool
}

func NewPostgreStore() (*PostgreStore, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	time.Sleep(10 * time.Second)
	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &PostgreStore{db: db}, nil
}

func (s *PostgreStore) CreatePost(post *models.Post) error {
	err := s.db.QueryRow(context.Background(), `
	  INSERT INTO posts (title, body, is_disabled_comments, user_id)
	  VALUES ($1, $2, $3, $4)
	  RETURNING id
	`, post.Title, post.Body, post.IsDisabledComments, post.UserID).Scan(&post.ID)
	return err
}

func (s *PostgreStore) GetAllPosts() ([]*models.Post, error) {
	rows, err := s.db.Query(context.Background(), `SELECT id, title, body, is_disabled_comments, user_id FROM posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.IsDisabledComments, &post.UserID); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}
	return posts, nil
}

func (s *PostgreStore) GetPostByID(id string) (*models.Post, error) {
	var post models.Post
	err := s.db.QueryRow(context.Background(), `SELECT id, title, body, is_disabled_comments, user_id FROM posts WHERE id=$1`, id).Scan(&post.ID, &post.Title, &post.Body, &post.IsDisabledComments, &post.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (s *PostgreStore) UpdatePost(post *models.Post) error {
	_, err := s.db.Exec(context.Background(), `
	  UPDATE posts SET title=$2, body=$3, is_disabled_comments=$4, user_id=$5 WHERE id=$1
	`, post.ID, post.Title, post.Body, post.IsDisabledComments, post.UserID)
	return err
}

func (s *PostgreStore) CreateComment(comment *models.Comment) error {
	err := s.db.QueryRow(context.Background(), `
	  INSERT INTO comments (post_id, parent_id, body, user_id)
	  VALUES ($1, $2, $3, $4)
	  RETURNING id
	`, comment.PostID, comment.ParentID, comment.Body, comment.UserID).Scan(&comment.ID)
	return err
}

func (s *PostgreStore) GetCommentsByPostIDWithPagination(postID string, limit, offset int) ([]*models.Comment, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT id, post_id, parent_id, body, user_id 
	FROM comments 
	WHERE post_id = $1 AND parent_id IS NULL 
	ORDER BY id ASC 
	LIMIT $2 OFFSET $3
	`, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rootComments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Body, &comment.UserID); err != nil {
			return nil, err
		}
		rootComments = append(rootComments, &comment)
	}

	var paginatedComments []*models.Comment
	for _, rootComment := range rootComments {
		paginatedComments = append(paginatedComments, rootComment)
		nestedComments, err := s.getNestedComments(rootComment.ID)
		if err != nil {
			return nil, err
		}
		paginatedComments = append(paginatedComments, nestedComments...)
	}

	return models.BuildCommentTree(paginatedComments), nil
}

func (s *PostgreStore) getNestedComments(parentID string) ([]*models.Comment, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT id, post_id, parent_id, body, user_id 
	FROM comments 
	WHERE parent_id = $1
	ORDER BY id ASC
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nestedComments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Body, &comment.UserID); err != nil {
			return nil, err
		}
		nestedComments = append(nestedComments, &comment)
		moreNestedComments, err := s.getNestedComments(comment.ID)
		if err != nil {
			return nil, err
		}
		nestedComments = append(nestedComments, moreNestedComments...)
	}
	return nestedComments, nil
}
