package model

type Post struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	CommentsDisabled bool       `json:"commentsDisabled"`
	UserID           string     `json:"userID"`
	Comments         []*Comment `json:"comments"`
}

type Comment struct {
	ID       string     `json:"id"`
	PostID   string     `json:"postID"`
	ParentID *string    `json:"parentID,omitempty"`
	Content  string     `json:"content"`
	UserID   string     `json:"userID"`
	Children []*Comment `json:"children"`
}
