package models

import (
	"time"
)

type Article struct {
	Id          int64     `db:"article_id"`
	Title       string    `db:"article_title"`
	Slug        string    `db:"article_slug"`
	Published   bool      `db:"article_published"`
	Content     string    `db:"article_content"`
	Author      int       `db:"article_author"`
	CreatedAt   time.Time `db:"article_created_at"`
	UpdatedAt   time.Time `db:"article_updated_at"`
	PublishedAt time.Time `db:"article_published_at"`
}
