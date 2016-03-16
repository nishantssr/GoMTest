package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"bitbucket.org/evard/evardbugs/app/models"
	_ "github.com/lib/pq"
	"github.com/ottob/go-semver/semver"
	"github.com/ottob/gorp"
	"github.com/revel/revel"
	"github.com/revel/modules/db/app"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}

	Dbm.AddTable(models.Case{}).SetKeys(true, "ID")

	Dbm.TypeConverter = typeConverter{}
	Dbm.TraceOn("[gorp]", revel.INFO)
	Dbm.CreateTables()
}

type typeConverter struct{}

func (me typeConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case *semver.Version:
		return t.String(), nil
	}

	return val, nil
}

func (me typeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case **semver.Version:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("FromDb: Unable to convert version to *string")
			}
			version, ok := target.(**semver.Version)
			if !ok {
				return errors.New(fmt.Sprint("FromDb: Unable to convert target to semver.Version: ", reflect.TypeOf(target)))
			}
			ver, err := semver.NewVersion(*s)
			if err != nil {
				return errors.New(fmt.Sprint("FromDb: Unable to create target semver.Version: ", reflect.TypeOf(target)))
			}
			*version = ver
			return nil
		}
		return gorp.CustomScanner{new(string), target, binder}, true
	}

	return gorp.CustomScanner{}, false
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
