package models

import (
  "github.com/revel/revel"
)

type Category struct {
  Id    int64   `db:"category_id"`
  Name  string  `db:"category_name"`
  Slug  string  `db:"category_slug"`
}

func (c *Category) Validate(v *revel.Validation) {
  v.Check(c.Name,
    revel.Required{},
    revel.MinSize{2},
    revel.MaxSize{100},
  )
}