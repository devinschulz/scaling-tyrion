package models

import (
	"database/sql"
	"fmt"
	"github.com/revel/revel"
	"regexp"
)

type User struct {
	Id             int64
	Name           string
	Email          string
	Password       sql.NullString
	HashedPassword []byte
}

func (user *User) String() string {
	return fmt.Sprintf("User(%s)", user.Email)
}

var emailPattern = regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Email,
		revel.Required{},
		revel.MaxSize{30},
		revel.MinSize{4},
		revel.Match{emailPattern},
	)

	ValidatePassword(v, user.Password.String).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{100},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
	)
}
