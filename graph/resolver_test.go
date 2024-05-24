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

// package graph

// import (
// 	"InternProj/graph/generated"
// 	"InternProj/internal/storages/memory"
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/stretchr/testify/assert"
// )

// // CreateTestServer sets up a test GraphQL server
// func CreateTestServer() *handler.Server {
// 	store := memory.NewMemoryStore()
// 	resolver := &Resolver{Store: store}

// 	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

// 	return srv
// }

// func TestGetAllPosts(t *testing.T) {
// 	server := CreateTestServer()

// 	// Create a new HTTP request
// 	query := `{
// 		getAllPosts {
// 			id
// 			title
// 			content
// 			comments {
// 				id
// 				content
// 			}
// 		}
// 	}`

// 	req, err := http.NewRequest("POST", "/query", bytes.NewBufferString(`{"query":"`+query+`"}`))
// 	if err != nil {
// 		t.Fatalf("could not create request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	// Create a new HTTP recorder
// 	rr := httptest.NewRecorder()

// 	// Serve the HTTP request
// 	server.ServeHTTP(rr, req)

// 	// Check the status code
// 	assert.Equal(t, http.StatusOK, rr.Code, "expected status OK")

// 	// Parse the response
// 	var resp struct {
// 		Data struct {
// 			GetAllPosts []struct {
// 				ID       string `json:"id"`
// 				Title    string `json:"title"`
// 				Content  string `json:"content"`
// 				Comments []struct {
// 					ID      string `json:"id"`
// 					Content string `json:"content"`
// 				} `json:"comments"`
// 			} `json:"getAllPosts"`
// 		} `json:"data"`
// 	}

// 	err = json.NewDecoder(rr.Body).Decode(&resp)
// 	if err != nil {
// 		t.Fatalf("could not decode response: %v", err)
// 	}

// 	// Check the response data
// 	assert.Len(t, resp.Data.GetAllPosts, 0, "expected no posts initially")
// }

// func TestCreatePostAndGetPostByID(t *testing.T) {
// 	server := CreateTestServer()

// 	// Create a post
// 	mutation := `mutation {
// 		createPost(input: {title: "Post 1", content: "Content 1", userID: "user1"}) {
// 			id
// 			title
// 			content
// 		}
// 	}`

// 	req, err := http.NewRequest("POST", "/query", bytes.NewBufferString(`{"query":"`+mutation+`"}`))
// 	if err != nil {
// 		t.Fatalf("could not create request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()
// 	server.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code, "expected status OK")

// 	var createResp struct {
// 		Data struct {
// 			CreatePost struct {
// 				ID      string `json:"id"`
// 				Title   string `json:"title"`
// 				Content string `json:"content"`
// 			} `json:"createPost"`
// 		} `json:"data"`
// 	}

// 	err = json.NewDecoder(rr.Body).Decode(&createResp)
// 	if err != nil {
// 		t.Fatalf("could not decode response: %v", err)
// 	}

// 	postID := createResp.Data.CreatePost.ID
// 	assert.NotEmpty(t, postID, "expected post ID to be set")

// 	// Get the post by ID
// 	query := `{
// 		getPostByID(id: "` + postID + `") {
// 			id
// 			title
// 			content
// 			comments {
// 				id
// 				content
// 			}
// 		}
// 	}`

// 	req, err = http.NewRequest("POST", "/query", bytes.NewBufferString(`{"query":"`+query+`"}`))
// 	if err != nil {
// 		t.Fatalf("could not create request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	rr = httptest.NewRecorder()
// 	server.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code, "expected status OK")

// 	var getResp struct {
// 		Data struct {
// 			GetPostByID struct {
// 				ID       string `json:"id"`
// 				Title    string `json:"title"`
// 				Content  string `json:"content"`
// 				Comments []struct {
// 					ID      string `json:"id"`
// 					Content string `json:"content"`
// 				} `json:"comments"`
// 			} `json:"getPostByID"`
// 		} `json:"data"`
// 	}

// 	err = json.NewDecoder(rr.Body).Decode(&getResp)
// 	if err != nil {
// 		t.Fatalf("could not decode response: %v", err)
// 	}

// 	assert.Equal(t, postID, getResp.Data.GetPostByID.ID, "expected post ID to match")
// 	assert.Equal(t, "Post 1", getResp.Data.GetPostByID.Title, "expected post title to match")
// 	assert.Equal(t, "Content 1", getResp.Data.GetPostByID.Content, "expected post content to match")
// 	assert.Len(t, getResp.Data.GetPostByID.Comments, 0, "expected no comments initially")
// }

// package graph

// import (
// 	"InternProj/internal/models"
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// // Mock DataStore
// type MockDataStore struct {
// 	Posts    []*models.Post
// 	Comments []*models.Comment
// }

// func (m *MockDataStore) GetAllPosts() ([]*models.Post, error) {
// 	return m.Posts, nil
// }

// func (m *MockDataStore) GetPostByID(id string) (*models.Post, error) {
// 	for _, post := range m.Posts {
// 		if post.ID == id {
// 			return post, nil
// 		}
// 	}
// 	return nil, nil
// }

// func (m *MockDataStore) GetCommentsByPostIDWithPagination(postID string, limit, offset int) ([]*models.Comment, error) {
// 	var comments []*models.Comment
// 	for _, comment := range m.Comments {
// 		if comment.PostID == postID {
// 			comments = append(comments, comment)
// 		}
// 	}
// 	start := offset
// 	end := offset + limit
// 	if start > len(comments) {
// 		start = len(comments)
// 	}
// 	if end > len(comments) {
// 		end = len(comments)
// 	}
// 	return comments[start:end], nil
// }

// func TestGetAllPosts(t *testing.T) {
// 	mockData := &MockDataStore{
// 		Posts: []*models.Post{
// 			{ID: "1", Title: "Post 1", Body: "Content 1", IsDisabledComments: false},
// 			{ID: "2", Title: "Post 2", Body: "Content 2", IsDisabledComments: false},
// 		},
// 		Comments: []*models.Comment{
// 			{ID: "1", PostID: "1", Body: "Comment 1", ParentID: nil},
// 			{ID: "2", PostID: "1", Body: "Comment 2", ParentID: *"1"},
// 		},
// 	}

// 	resolver := &Resolver{Store: mockData}
// 	query := &queryResolver{resolver}

// 	limit := 2
// 	offset := 0
// 	posts, err := query.GetAllPosts(context.Background(), &limit, &offset)
// 	assert.NoError(t, err)
// 	assert.Len(t, posts, 2)
// 	assert.Len(t, posts[0].Comments, 2)
// 	assert.Equal(t, "Post 1", posts[0].Title)
// 	assert.Equal(t, "Comment 1", posts[0].Comments[0].Body)
// }

// func TestGetPostByID(t *testing.T) {
// 	mockData := &MockDataStore{
// 		Posts: []*models.Post{
// 			{ID: "1", Title: "Post 1", Body: "Content 1", IsDisabledComments: false},
// 		},
// 		Comments: []*models.Comment{
// 			{ID: "1", PostID: "1", Body: "Comment 1", ParentID: nil},
// 			{ID: "2", PostID: "1", Body: "Comment 2", ParentID: "1"},
// 		},
// 	}

// 	resolver := &Resolver{Store: mockData}
// 	query := &queryResolver{resolver}

// 	limit := 2
// 	offset := 0
// 	post, err := query.GetPostByID(context.Background(), "1", &limit, &offset)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, post)
// 	assert.Equal(t, "Post 1", post.Title)
// 	assert.Len(t, post.Comments, 2)
// 	assert.Equal(t, "Comment 1", post.Comments[0].Body)
// }
