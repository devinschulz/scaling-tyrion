package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/idevschulz/portfolio/app/models"
	"github.com/revel/revel"
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
	auth := c.Auth()
	return c.Render(auth)
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

func (c App) getUserById(id int) *models.User {
	user, err := c.Txn.Select(models.User{}, `SELECT * FROM users WHERE id = $1`, id)
	checkErr(err, "Select User by ID Failed: ")
	if len(user) == 0 {
		return nil
	}
	return user[0].(*models.User)
}

func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
	user.Validate(c.Validation)
	user.ValidatePassword(c.Validation, user.Password)
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).Message("Passwords don't match")

	UserEmailCheck, err := c.Txn.Select(models.User{}, `SELECT * FROM users WHERE email = $1`, user.Email)
	checkErr(err, "Failed to get email: ")

	if UserEmailCheck != nil {
		c.Validation.Error("Email already taken").Key("user.Email")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Register)
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Register)
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err = c.Txn.Insert(&user)
	checkErr(err, "Saving User failed:")

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome " + user.Name)
	return c.Redirect(App.Index)
}

func (c App) Login(email, password string, remember bool) revel.Result {
	if email == "" || password == "" {
		c.Flash.Error("Missing fields")
		return c.Redirect(App.Index)
	}
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
			if user.Admin == true {
				c.Session["admin"] = email
				c.Flash.Success("Welcome Back " + user.Name + ". You are an Admin!")
			} else {
				c.Flash.Success("Welcome Back " + user.Name)
			}
			return c.Redirect(App.Index)
		} else {
			c.Flash.Error("Incorrect Password")
			return c.Redirect(App.Index)
		}
	}
	c.Flash.Error("Email Does Not Exist")
	return c.Redirect(App.Index)
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.Flash.Success("Logged out successfully")
	return c.Redirect(App.Index)
}

func (c App) Settings() revel.Result {
	auth := c.Auth()
	if email, ok := c.Session["user"]; ok {
		user := c.getUser(email)
		return c.Render(auth, user)
	}

	c.Flash.Error("You Must Be Logged In To Edit Your Profile")
	return c.Redirect(App.Index)
}

func (c App) UpdateSettings(user *models.User, verifyPassword string) revel.Result {

	s := c.connected()

	user.Id = s.Id
	user.Validate(c.Validation)

	// Only validate password if password is given
	if user.Password != "" || verifyPassword != "" {
		user.ValidatePassword(c.Validation, user.Password)
		c.Validation.Required(verifyPassword)
		c.Validation.Required(verifyPassword == user.Password).Message("Passwords don't match")
		bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.HashedPassword = bcryptPassword
	} else {
		existingPassword := s.HashedPassword
		if existingPassword != nil {
			user.HashedPassword = existingPassword
		}
	}

	if user.Email != s.Email {
		UB, err := c.Txn.Select(models.User{}, `SELECT * FROM users WHERE email = $1`, user.Email)
		checkErr(err, "Failed to get email: ")
		if UB != nil && len(UB) > 0 {
			c.Validation.Error("Email already taken").Key("user.Email")
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(App.Settings)
		}
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Please correct the errors below.")
		return c.Redirect(App.Settings)
	}

	// Update User
	_, err := c.Txn.Update(user)
	checkErr(err, "Failed to Update User: ")

	// Refresh the session in case the email address was changed.
	c.Session["user"] = user.Email

	c.Flash.Success("Settings Updated")
	return c.Redirect(App.Settings)

}
