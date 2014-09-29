package models

import (
	"github.com/revel/revel"
  "github.com/coopernurse/gorp"
  "github.com/russross/blackfriday"
	"time"
  "html/template"
)

const trimLength = 300

type Article struct {
	Id          int64     `db:"article_id"`
	Title       string    `db:"article_title"`
	Slug        string    `db:"article_slug"`
	Published   bool      `db:"article_published"`
	Content     string    `db:"article_content"`
	AuthorId    int64     `db:"article_author_id"`
  Categories  string    `db:"article_categories"`
  Tags        string    `db:"article_tags"`
	CreatedAt   time.Time `db:"article_created_at"`
	UpdatedAt   time.Time `db:"article_updated_at"`
	PublishedAt time.Time `db:"article_published_at"`
  Meta        map[string]interface{} `-`
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

func (a *Article) PreInsert(s gorp.SqlExecutor) error {
    a.CreatedAt = time.Now()
    a.UpdatedAt = a.CreatedAt
    return nil
}

func (a *Article) PreUpdate(s gorp.SqlExecutor) error {
    a.UpdatedAt = time.Now()
    return nil
}

func (article *Article) addMeta() {
  if article.Meta == nil {
    article.Meta = make(map[string]interface{})
  }
  //article.Meta["categories"] = strings.Split(article.Categories, ",")
  //article.Meta["tags"] = strings.Split(article.Tags, ",")
  if len(article.Content) > trimLength {
    article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Content[0:trimLength]))))
  } else {
    article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Content[0:len(article.Content)]))))
  }
}