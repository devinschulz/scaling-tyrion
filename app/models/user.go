package models

import (
	"fmt"
	"github.com/revel/revel"
	"regexp"
)

type User struct {
	Id             int64
	Name           string
	Email          string
	Password       string
	HashedPassword []byte
	Admin          bool
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

	v.Email(user.Email)

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{100},
	)
}

func (user *User) ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
	)
}
