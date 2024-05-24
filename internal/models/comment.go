package models

type Comment struct {
	ID            string     `json:"id"`
	PostID        string     `json:"postID"`
	ParentID      *string    `json:"parentID,omitempty"`
	Body          string     `json:"body"`
	UserID        string     `json:"userID"`
	ChildComments []*Comment `json:"childComments"`
}

func BuildCommentTree(comments []*Comment) []*Comment {
	commentMap := make(map[string]*Comment)
	var rootComments []*Comment

	for _, comment := range comments {
		commentMap[comment.ID] = comment
		comment.ChildComments = []*Comment{}
	}

	for _, comment := range comments {
		if comment.ParentID == nil {
			rootComments = append(rootComments, comment)
		} else {
			if parent, ok := commentMap[*comment.ParentID]; ok {
				parent.ChildComments = append(parent.ChildComments, comment)
			}
		}
	}

	return rootComments
}
