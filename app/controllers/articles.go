package controllers

import (
	"github.com/idevschulz/portfolio/app/models"
	"github.com/revel/revel"
	// "log"
)

type Articles struct {
	App
}

func (c Articles) Index() revel.Result {
	auth := c.Auth()
	articles, err := c.Txn.Select(models.Article{}, `SELECT * FROM articles`)
	checkErr(err, "Failed to select articles")
	return c.Render(auth, articles)
}

func (c Articles) New() revel.Result {
	auth := c.Auth()
	action := "/articles/new"
	actionButton := "Create Article"
	return c.Render(auth, action, actionButton)
}

func (c Articles) CreateArticle(article models.Article) revel.Result {
	article.Validate(c.Validation)

	// TODO: Re-factor this
	slug, err := c.Txn.Select(models.Article{}, `SELECT * FROM articles WHERE article_slug = $1`, article.Slug)
	checkErr(err, "Failed to select article slug")
	if slug != nil && len(slug) > 0 {
		c.Validation.Error("Slug Already in Use").Key("article.Slug")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	// TODO: Re-factor this
	title, err := c.Txn.Select(models.Article{}, `SELECT * FROM articles WHERE article_title = $1`, article.Title)
	checkErr(err, "Failed to select article slug")
	if title != nil && len(title) > 0 {
		c.Validation.Error("Title Already in Use").Key("article.Title")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Check For Errors")
		return c.Redirect(Articles.New)
	}

	s := c.connected()
	article.AuthorId = s.Id

	err = c.Txn.Insert(&article)
	checkErr(err, "Saving Article failed: ")

	c.Flash.Success("Article Created: " + article.Title)
	return c.Redirect(Articles.Index)
}

func (c Articles) CheckIfSlugExists(slug string) bool {
	s, err := c.Txn.Select(models.Article{}, `SELECT * FROM articles WHERE article_slug = $1`, slug)
	checkErr(err, "Failed to select article slug")
	if s != nil && len(s) > 0 {
		return true
	}
	return false
}
