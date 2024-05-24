package graph

import (
	"InternProj/graph/generated"
	"InternProj/internal/storages/memory"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	resolver := &Resolver{
		Store: memory.NewMemoryStore(),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		CreatePost struct {
			ID                 string
			Title              string
			Body               string
			UserID             string
			IsDisabledComments bool
		}
	}

	c.MustPost(`mutation {
    createPost(title: "Test Post", body: "This is a test post", userID: "1") {
      id
      title
      body
	  userID
	  isDisabledComments
    }
  }`, &resp)

	assert.Equal(t, "Test Post", resp.CreatePost.Title)
	assert.Equal(t, "This is a test post", resp.CreatePost.Body)
	assert.Equal(t, "1", resp.CreatePost.UserID)
	assert.False(t, resp.CreatePost.IsDisabledComments)
}

func TestCreateComment(t *testing.T) {
	resolver := &Resolver{
		Store: memory.NewMemoryStore(),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var postResp struct {
		CreatePost struct {
			ID                 string
			Title              string
			Body               string
			UserID             string
			IsDisabledComments bool
		}
	}

	c.MustPost(`mutation {
		createPost(title: "Test Post", body: "This is a test post", userID: "1") {
		  id
		  title
		  body
		  userID
		  isDisabledComments
		}
	  }`, &postResp)

	var commentResp struct {
		CreateComment struct {
			ID       string
			PostID   string
			ParentID *string
			Body     string
		}
	}

	c.MustPost(`mutation {
    createComment(postID: "`+postResp.CreatePost.ID+`", body: "This is a test comment") {
      id
      postID
      parentID
      body
    }
  }`, &commentResp)

	assert.Equal(t, postResp.CreatePost.ID, commentResp.CreateComment.PostID)
	assert.Nil(t, commentResp.CreateComment.ParentID)
	assert.Equal(t, "This is a test comment", commentResp.CreateComment.Body)
}

func TestCannotCreateComment(t *testing.T) {
	resolver := &Resolver{
		Store: memory.NewMemoryStore(),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var postResp struct {
		CreatePost struct {
			ID                 string
			Title              string
			Body               string
			UserID             string
			IsDisabledComments bool
		}
	}

	c.MustPost(`mutation {
		createPost(title: "Test Post", body: "This is a test post", userID: "1") {
		  id
		  title
		  body
		  userID
		  isDisabledComments
		}
	  }`, &postResp)

	var changeResp struct {
		ChangePostCommentsAccess struct {
			ID                 string
			Title              string
			Body               string
			UserID             string
			IsDisabledComments bool
		}
	}

	c.MustPost(`mutation {
		changePostCommentsAccess(postID: "`+postResp.CreatePost.ID+`", userID: "1", isDisabledComments: true) {
		  id
		  title
		  body
		  userID
		  isDisabledComments
		}
	}`, &changeResp)

	var commentResp struct {
		CreateComment struct {
			ID       string
			PostID   string
			ParentID *string
			Body     string
		}
	}

	err := c.Post(`mutation {
		createComment(postID: "`+postResp.CreatePost.ID+`", body: "This is a test comment") {
		  id
		  postID
		  parentID
		  body
		}
	}`, &commentResp)

	assert.Error(t, err)

	assert.Contains(t, err.Error(), "the post is closed for comment")
}
