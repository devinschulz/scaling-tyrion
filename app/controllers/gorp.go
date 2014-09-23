package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/idevschulz/portfolio/app/models"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/db/app"
	"log"
)

var (
	Dbm *gorp.DbMap
)

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Count not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTable(models.User{}).SetKeys(true, "Id")
	t.ColMap("Password").Transient = true
	setColumnSizes(t, map[string]int{
		"Email": 50,
		"Name":  50,
	})

	Dbm.TraceOn("[gorp]", revel.INFO)

	err := Dbm.CreateTablesIfNotExists()
	checkErr(err, "Create Tables failed")

	// Generate Demo User
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte("demo"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "demo", "demo@demo.com", "demo", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		checkErr(err, "Failed to insert user")
	}

}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
