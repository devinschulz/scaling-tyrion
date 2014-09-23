package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/idevschulz/portfolio/app/models"
	"github.com/idevschulz/portfolio/app/routes"
	"github.com/revel/revel"
	// "log"
)

type App struct {
	*revel.Controller
	GorpController
}

func (c App) Auth() bool {
	auth := false
	if user := c.connected(); user != nil {
		auth = true
	}
	return auth
}

func (c App) Index() revel.Result {
	auth := c.Auth()
	return c.Render(auth)
}

func (c App) Register() revel.Result {
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
	users, err := c.Txn.Select(models.User{}, `SELECT * FROM users WHERE email = $1`, email)
	checkErr(err, "Select failed:")
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).Message("Passwords don't match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	checkErr(err, "Saving User failed:")

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
			c.Flash.Success("Welcome Back " + user.Name)
			return c.Redirect(routes.App.Index())
		} else {
			c.Flash.Error("Incorrect Password")
			return c.Redirect(routes.App.Index())
		}
	}
	c.Flash.Error("User Not found")
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.Flash.Success("Logged out successfully")
	return c.Redirect(routes.App.Index())
}
