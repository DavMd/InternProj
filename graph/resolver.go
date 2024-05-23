package graph

import (
	"InternProj/graph/generated"
	"InternProj/graph/model"
	"context"
	"fmt"
	"math/rand"
	"strconv"
)

type Resolver struct {
	posts []*model.Post
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, title string, body string, userID string) (*model.Post, error) {
	post := &model.Post{
		ID:                 strconv.Itoa(rand.Int()),
		Title:              title,
		Body:               body,
		IsDisabledComments: false,
		UserID:             userID,
		Comments:           []*model.Comment{},
	}
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *mutationResolver) ChangePostCommentsAccess(ctx context.Context, postID string, userID string, commentsDisabled bool) (*model.Post, error) {
	for _, post := range r.posts {
		if post.ID == postID {
			if post.UserID == userID {
				post.IsDisabledComments = commentsDisabled
				return post, nil
			}
			return nil, fmt.Errorf("wrong user")
		}
	}
	return nil, fmt.Errorf("post not found")
}
