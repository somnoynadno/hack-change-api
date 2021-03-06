package entities

import "hack-change-api/models/auxiliary"

type BlogPost struct {
	auxiliary.BaseModel
	Text        string                 `json:"text"`
	Instruments []*FinancialInstrument `json:"instruments" gorm:"many2many:post_instruments;"`
	AuthorID    uint                   `json:"authorID"`
	Author      *User                  `json:"author,omitempty"`
	Comments    []*Comment             `json:"comments,omitempty"`
	Likes       []*User                `json:"likes,omitempty" gorm:"many2many:like_blog_posts"`
}
