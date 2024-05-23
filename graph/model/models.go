package model

type Post struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Body               string     `json:"body"`
	IsDisabledComments bool       `json:"isDisabledComments"`
	UserID             string     `json:"userID"`
	Comments           []*Comment `json:"comments"`
}

type Comment struct {
	ID            string     `json:"id"`
	PostID        string     `json:"postID"`
	ParentID      *string    `json:"parentID,omitempty"`
	Body          string     `json:"body"`
	UserID        string     `json:"userID"`
	ChildComments []*Comment `json:"childComments"`
}
