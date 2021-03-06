package controllers

import (
	"database/sql"
	"restgo/app/models"

	"github.com/coopernurse/gorp"
	"github.com/revel/revel"
)

//
var (
	Dbm *gorp.DbMap
)

// GorpController will be extended for others
type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

// InitDb is the function called in the start of the application
var InitDb = func() {
	connectionString := getConnectionString()
	if db, err := sql.Open("mysql", connectionString); err != nil {
		revel.ERROR.Fatal(err)
	} else {
		Dbm = &gorp.DbMap{
			Db:      db,
			Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	}
	// Defines the table for use by GORP
	// This is a function we will create soon.
	defineBidItemTable(Dbm)
	if err := Dbm.CreateTablesIfNotExists(); err != nil {
		revel.ERROR.Fatal(err)
	}
}

// Begin is a method that starts a transaction
func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

// Commit Method that commits something
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

// Rollback method that returns a transaction
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

func defineBidItemTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTable(models.BidItem{}).SetKeys(true, "id")
	// e.g. VARCHAR(25)
	t.ColMap("name").SetMaxSize(25)
}
