package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/idevschulz/portfolio/app/models"
	"github.com/idevschulz/portfolio/app/routes"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	GorpController
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.RenderArgs["user"] = user
	}
	return nil
}

func (c App) connected() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}
	if email, ok := c.Session["user"]; ok {
		return c.getUser(email)
	}
	return nil
}

func (c App) getUser(email string) *models.User {
	users, err := c.Txn.Select(models.User{}, `SELECT * FROM users where Email = $1`, email)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

func (c App) Register() revel.Result {
	return c.Render()
}

func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).Message("Passwords don't match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		// return c.Redirect(routes.App.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome " + user.Name)
	return c.Redirect(routes.App.Index())
}

func (c App) Login(email, password string, remember bool) revel.Result {
	user := c.getUser(email)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = email
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome back" + user.Name)
			return c.Redirect(routes.App.Index())
		}
	}
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}
