package controllers

import (
  "github.com/revel/revel"
  "github.com/idevschulz/portfolio/app/models"
  "fmt"
  "strings"
  "unicode"
)

type Categories struct {
  App
}

func (c Categories) Index() revel.Result {
  auth := c.Auth()
  if auth {
    categories, err := c.Txn.Select(models.Category{}, `SELECT * FROM categories ORDER BY category_name ASC`)
    checkERROR(err)
    return c.Render(auth, categories)  
  }
  c.Flash.Error("You must be be logged in to delete a category")
  return c.Redirect(App.Index)
}

func (c Categories) Edit(id int64) revel.Result {
  auth := c.Auth()
  if auth {
    var category models.Category
    err := c.Txn.SelectOne(&category, `SELECT * FROM categories WHERE category_id=$1`, id)
    checkERROR(err)
    return c.Render(auth, category)  
  }
  c.Flash.Error("You must be be logged in to edit a category")
  return c.Redirect(App.Index) 
}

func (c Categories) New(category models.Category) revel.Result {
  auth := c.Auth()
  if auth {
    category.Validate(c.Validation)

    if c.Validation.HasErrors() {
      c.Validation.Keep()
      c.FlashParams()
      c.Flash.Error("Check For Errors")
      return c.Redirect(Categories.Index)
    }

    // Check if unique
    name, err := c.Txn.Select(models.Category{}, `SELECT * FROM categories WHERE category_name=$1`, category.Name)
    checkERROR(err)
    if name != nil && len(name) > 0 {
      c.Validation.Error("Category Name already in use").Key("category.Name")
      c.Validation.Keep()
      c.FlashParams()
      return c.Redirect(Categories.Index)
    }

    category.Slug = GenerateSlug(category.Name)

    err = c.Txn.Insert(&category)
    checkERROR(err)

    c.Flash.Success("Category Created: " + category.Name)
    return c.Redirect(Categories.Index)
  } 
  c.Flash.Error("You must be be logged in to create a category")
  return c.Redirect(App.Index)
}

func (c Categories) Delete(id int64) revel.Result {
  auth := c.Auth()
  if auth {
    _, err := c.Txn.Delete(&models.Category{Id: id})
    checkERROR(err)
    c.Flash.Success(fmt.Sprintf("Category has been deleted", id))
    return c.Redirect(Categories.Index)    
  }
  c.Flash.Error("You must be be logged in to delete a category")
  return c.Redirect(App.Index)
}

func (c Categories) Update(category *models.Category, id int64) revel.Result {
  auth := c.Auth()
  if auth {
    _, err := c.Txn.Exec(`UPDATE categories SET category_name=$1, category_slug=$2 WHERE category_id=$3`, category.Name, category.Slug, id)
    checkERROR(err)
    c.Flash.Success(fmt.Sprintf("Category has been updated", id))
    return c.Redirect(Categories.Index)
  }
  c.Flash.Error("You must be be logged in to update a category")
  return c.Redirect(App.Index)
}

func GenerateSlug(str string) (slug string) {
  return strings.Map(func(r rune) rune {
    switch {
    case r == ' ', r == '-':
      return '-'
    case r == '_', unicode.IsLetter(r), unicode.IsDigit(r):
      return r
    default:
      return -1
    }
    return -1
  }, strings.ToLower(strings.TrimSpace(str)))
} 