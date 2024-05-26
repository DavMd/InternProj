package graph

import (
	"InternProj/graph/generated"
	"InternProj/internal/models"
	"InternProj/internal/storages"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

type Resolver struct {
	Store       storages.Storage
	subscribers map[string][]chan *models.Comment
	mu          sync.Mutex
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() generated.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, title string, body string, userID string) (*models.Post, error) {
	post := &models.Post{
		ID:                 strconv.Itoa(rand.Int()),
		Title:              title,
		Body:               body,
		IsDisabledComments: false,
		UserID:             userID,
		Comments:           []*models.Comment{},
	}
	err := r.Store.CreatePost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, body string) (*models.Comment, error) {
	if len(body) > 2000 {
		return nil, fmt.Errorf("post length exceeds 2000 characters")
	}

	post, err := r.Store.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if post.IsDisabledComments {
		return nil, fmt.Errorf("the post is closed for comment")
	}

	comment := &models.Comment{
		ID:            strconv.Itoa(rand.Int()),
		PostID:        postID,
		ParentID:      parentID,
		Body:          body,
		UserID:        post.UserID,
		ChildComments: []*models.Comment{},
	}

	err = r.Store.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ch := range r.subscribers[postID] {
		ch <- comment
	}

	return comment, nil
}

func (r *mutationResolver) ChangePostCommentsAccess(ctx context.Context, postID string, userID string, commentsDisabled bool) (*models.Post, error) {
	post, err := r.Store.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if userID == post.UserID {
		post.IsDisabledComments = commentsDisabled

		err = r.Store.UpdatePost(post)
		if err != nil {
			return nil, err
		}

		return post, nil
	}
	return nil, fmt.Errorf("wrong user post")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	posts, err := r.Store.GetAllPosts()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *queryResolver) GetPostByID(ctx context.Context, id string, limit *int, offset *int) (*models.Post, error) {
	post, err := r.Store.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	l := 10
	o := 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}

	comments, err := r.Store.GetCommentsByPostIDWithPagination(id, l, o)
	if err != nil {
		return nil, err
	}

	post.Comments = comments
	return post, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	ch := make(chan *models.Comment, 1)

	r.mu.Lock()
	if r.subscribers == nil {
		r.subscribers = make(map[string][]chan *models.Comment)
	}
	r.subscribers[postID] = append(r.subscribers[postID], ch)
	r.mu.Unlock()

	go func() {
		<-ctx.Done()
	}()

	return ch, nil
}
