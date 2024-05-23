package models

type Comment struct {
	ID            string     `json:"id"`
	PostID        string     `json:"postID"`
	ParentID      *string    `json:"parentID,omitempty"`
	Body          string     `json:"body"`
	UserID        string     `json:"userID"`
	ChildComments []*Comment `json:"childComments"`
}
