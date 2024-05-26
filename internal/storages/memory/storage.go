package memory

import (
	"InternProj/internal/models"
	"errors"
	"sync"
)

type MemoryStore struct {
	posts    []*models.Post
	comments []*models.Comment
	mu       sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		posts:    []*models.Post{},
		comments: []*models.Comment{},
	}
}

func (s *MemoryStore) CreatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts = append(s.posts, post)
	return nil
}

func (s *MemoryStore) GetAllPosts() ([]*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.posts, nil
}

func (s *MemoryStore) GetPostByID(id string) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, post := range s.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (s *MemoryStore) UpdatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.posts {
		if p.ID == post.ID {
			s.posts[i] = post
			return nil
		}
	}
	return errors.New("post not found")
}

func (s *MemoryStore) CreateComment(comment *models.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments = append(s.comments, comment)
	return nil
}

func (s *MemoryStore) GetCommentsByPostIDWithPagination(postID string, limit, offset int) ([]*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comments := []*models.Comment{}
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	start := offset
	end := offset + limit
	if start > len(comments) {
		start = len(comments)
	}
	if end > len(comments) {
		end = len(comments)
	}

	paginatedComments := comments[start:end]
	return models.BuildCommentTree(paginatedComments), nil
}
