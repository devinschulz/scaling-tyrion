package models

import (
	"github.com/revel/revel"
	"time"
)

type Article struct {
	Id          int64     `db:"article_id"`
	Title       string    `db:"article_title"`
	Slug        string    `db:"article_slug"`
	Published   bool      `db:"article_published"`
	Content     string    `db:"article_content"`
	AuthorId    int64     `db:"article_author_id"`
	CreatedAt   time.Time `db:"article_created_at"`
	UpdatedAt   time.Time `db:"article_updated_at"`
	PublishedAt time.Time `db:"article_published_at"`
}

func (article *Article) Validate(v *revel.Validation) {
	v.Check(article.Title,
		revel.Required{},
		revel.MaxSize{100},
		revel.MinSize{2},
	)
	v.Check(article.Slug,
		revel.Required{},
		revel.MaxSize{100},
		revel.MinSize{2},
	)
	v.Check(article.Content,
		revel.Required{},
	)
}
