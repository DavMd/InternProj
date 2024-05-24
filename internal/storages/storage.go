package storages

import "InternProj/internal/models"

type Storage interface {
	CreatePost(post *models.Post) error
	GetAllPosts() ([]*models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	UpdatePost(post *models.Post) error

	CreateComment(comment *models.Comment) error
	GetCommentsByPostID(postID string) ([]*models.Comment, error)
	GetCommentsByPostIDWithPagination(postID string, limit, offset int) ([]*models.Comment, error)
}
