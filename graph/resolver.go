package graph

import (
	"InternProj/graph/generated"
	"InternProj/graph/model"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Resolver struct {
	posts       []*model.Post
	subscribers map[string][]chan *model.Comment
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

func (r *mutationResolver) CreatePost(ctx context.Context, title string, body string, userID string) (*model.Post, error) {
	post := &model.Post{
		ID:                 uuid.NewString(),
		Title:              title,
		Body:               body,
		IsDisabledComments: false,
		UserID:             userID,
		Comments:           []*model.Comment{},
	}
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, body string) (*model.Comment, error) {
	if len(body) > 2000 {
		return nil, fmt.Errorf("post length exceeds 2000 characters")
	}

	for _, post := range r.posts {
		if post.ID == postID {
			if post.IsDisabledComments {
				return nil, fmt.Errorf("the post is closed for comment")
			}

			comment := &model.Comment{
				ID:            uuid.NewString(),
				PostID:        postID,
				ParentID:      parentID,
				Body:          body,
				ChildComments: []*model.Comment{},
			}

			if parentID == nil {
				post.Comments = append(post.Comments, comment)
			} else {
				parentComment := findCommentByID(post.Comments, *parentID)
				if parentComment == nil {
					return nil, fmt.Errorf("parent comment not found")
				}
				parentComment.ChildComments = append(parentComment.ChildComments, comment)
			}

			r.mu.Lock()
			defer r.mu.Unlock()
			for _, ch := range r.subscribers[postID] {
				ch <- comment
			}

			return comment, nil
		}
	}
	return nil, fmt.Errorf("post not found")
}

func findCommentByID(comments []*model.Comment, id string) *model.Comment {
	for _, comment := range comments {
		if comment.ID == id {
			return comment
		}
		if nestedComment := findCommentByID(comment.ChildComments, id); nestedComment != nil {
			return nestedComment
		}
	}
	return nil
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

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	return r.posts, nil
}

func (r *queryResolver) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	for _, post := range r.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, fmt.Errorf("post not found")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	ch := make(chan *model.Comment, 1)

	r.mu.Lock()
	if r.subscribers == nil {
		r.subscribers = make(map[string][]chan *model.Comment)
	}
	r.subscribers[postID] = append(r.subscribers[postID], ch)
	r.mu.Unlock()

	go func() {
		<-ctx.Done()
	}()

	return ch, nil
}
