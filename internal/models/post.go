package models

type Post struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Body               string     `json:"body"`
	IsDisabledComments bool       `json:"isDisabledComments"`
	UserID             string     `json:"userID"`
	Comments           []*Comment `json:"comments"`
}
