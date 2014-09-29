package controllers

import (
	"github.com/idevschulz/portfolio/app/models"
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	"html/template"
	"time"
  "fmt"
  "log"
)

type Articles struct {
	App
}

func (c Articles) Index() revel.Result {
	auth := c.Auth()

	var (
		articles []interface{}
		err      error
	)

	if auth {
		articles, err = c.Txn.Select(models.Article{}, `SELECT * FROM articles ORDER BY article_published_at DESC`)
	} else {
		articles, err = c.Txn.Select(models.Article{}, `SELECT * FROM articles WHERE article_published = true ORDER By article_published_at DESC`)
	}

	checkErr(err, "Failed to select articles")
	return c.Render(auth, articles)
}

func (c Articles) New() revel.Result {
	auth := c.Auth()
  if auth {
    action := "/articles/new"
    actionButton := "Create Article"
    categories := getCategories(c)
    return c.Render(auth, action, actionButton, categories)  
  }
  c.Flash.Error("You must be logged in to create new articles")
	return c.Redirect(Articles.Index)
}

func (c Articles) CreateArticle(article models.Article) revel.Result {
	article.Validate(c.Validation)

	// TODO: Re-factor this
	slug, err := c.Txn.Select(models.Article{}, `SELECT * FROM articles WHERE article_slug = $1`, article.Slug)
	checkERROR(err)
	if slug != nil && len(slug) > 0 {
		c.Validation.Error("Slug Already in Use").Key("article.Slug")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	title, err := c.Txn.Select(models.Article{}, `SELECT * FROM articles WHERE article_title = $1`, article.Title)
	checkERROR(err)
	if title != nil && len(title) > 0 {
		c.Validation.Error("Title Already in Use").Key("article.Title")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}
	// end

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Check For Errors")
		return c.Redirect(Articles.New)
	}

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	s := c.connected()
	article.AuthorId = s.Id

	err = c.Txn.Insert(&article)
	checkERROR(err)

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

func (c Articles) Show(slug string) revel.Result {
	if slug != "" {
		auth := c.Auth()
		article := c.GetArticleBySlug(slug)
		content := template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Content))))
		return c.Render(article, content, auth)
	}

	return c.NotFound("Invalid Article")
}

func (c Articles) Edit(id int64) revel.Result {

	auth := c.Auth()
  if auth {
  	action := fmt.Sprintf("/articles/update/%v", id)
  	actionButton := "Update Article"
  	article := c.GetArticleById(id)
    categories := getCategories(c)

    var (
      publishAction string
      publishButton string
      class string
    )

    if article.Published {
      publishAction = fmt.Sprintf("/articles/unpublish/%v", id)
      publishButton = "UnPublish"
      class = "btn-warning"
    } else {
      publishAction = fmt.Sprintf("/articles/publish/%v", id)
      publishButton = "Publish"
      class = "btn-success"
    }

  	return c.Render(auth, action, actionButton, article, publishAction, publishButton, class, categories)
  } 
  c.Flash.Error("You must be logged in to edit this post")
  return c.Redirect(Articles.Index)
}

func (c Articles) Update(article *models.Article, id int64) revel.Result {

  route := fmt.Sprintf("/articles/%v", article.Slug)
  auth := c.Auth()
  
  if auth {
    article.Validate(c.Validation)
    if c.Validation.HasErrors() {
      c.Validation.Keep()
      c.FlashParams()
      c.Flash.Error("Please correct the errors below.")
      return c.Redirect(App.Settings)
    }

    _, err := c.Txn.Exec(`UPDATE articles SET article_title=$1, article_slug=$2, article_content=$3 WHERE article_id=$4`, article.Title, article.Slug, article.Content, article.Id )
    checkERROR(err)

    c.Flash.Success("Article Updated")
    return c.Redirect(route)
  }   

  c.Flash.Error("You must be logged in to Update an article")
  return c.Redirect(route)
}

func (c Articles) Delete(id int64) revel.Result {
  auth := c.Auth()
  if auth {
    _, err := c.Txn.Delete(&models.Article{Id: id})
    checkERROR(err)
    c.Flash.Success(fmt.Sprintf("Article has been deleted", id))
    return c.Redirect(Articles.Index)    
  }

  c.Flash.Error("You do not have permission to delete this post")
  return c.Redirect(Articles.Index)
}

func (c Articles) Publish(id int64) revel.Result {
  auth := c.Auth()
  if auth {
    obj, err :=c.Txn.Get(models.Article{}, id)
    article := obj.(*models.Article)
    article.Published = true
    article.PublishedAt = time.Now()
    article.UpdatedAt = time.Now()

    _, err = c.Txn.Update(article)
    checkERROR(err)
    return c.Redirect(Articles.Index)
  }
  c.Flash.Error("You do not have permission to publish this post")
  return c.Redirect(Articles.Index)
}

func (c Articles) UnPublish(id int64) revel.Result {
  auth := c.Auth()
  if auth {
    obj, err :=c.Txn.Get(models.Article{}, id)
    article := obj.(*models.Article)
    article.Published = false
    article.UpdatedAt = time.Now()
    log.Print(time.Now())
    
    _, err = c.Txn.Update(article)
    checkERROR(err)
    return c.Redirect(Articles.Index)
  }

  c.Flash.Error("You do not have permission to publish this post")
  return c.Redirect(Articles.Index)
}

func (c Articles) GetArticleBySlug(slug string) models.Article {
	var article models.Article
	err := c.Txn.SelectOne(&article, `SELECT * FROM articles WHERE article_slug=$1`, slug)
  checkERROR(err)
	return article
}

func (c Articles) GetArticleById(id int64) models.Article {
	var article models.Article
	err := c.Txn.SelectOne(&article, `SELECT * FROM articles WHERE article_id=:id`, map[string]interface{}{
		"id": id,
	})
	checkERROR(err)
	return article
}

func getCategories(c Articles) []interface{} {
  categories, err := c.Txn.Select(models.Category{}, `SELECT * FROM categories ORDER BY category_name`)
  checkERROR(err)
  return categories
}