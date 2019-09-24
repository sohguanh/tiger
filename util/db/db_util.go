// dbUtil package has dependency on a third party package. please install it first before using this package.
//
// 	go get -v github.com/go-sql-driver/mysql
// this is a MySQL database driver. null column values need to be specially handled so sample code are referenced from
// https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
package dbUtil

import (
	"database/sql"
	"encoding/json"
	"errors"
	mysql "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"sync"
	"tiger/config"
	dateTimeUtil "tiger/util/datetime"
	logUtil "tiger/util/log"
)

const (
	JsonNull = "null"
)

var db *sql.DB
var onceDb sync.Once
var dbConfigErr bool = false

// NewDb return a singleton object for application use. connection pooling is built-in by the database driver.
func NewDb(c *config.Config) (*sql.DB, error) {
	onceDb.Do(func() { //singleton
		logUtil.DebugPrint("db first time init\n")
		dsn := c.Database.Username + ":" + c.Database.Password + "@tcp(" + c.Database.Host + ":" + strconv.Itoa(c.Database.Port) + ")/" + c.Database.Name + "?parseTime=true"
		logUtil.DebugPrint(dsn)
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
			dbConfigErr = true
		}
	})
	if !dbConfigErr {
		return db, nil
	} else {
		return nil, errors.New("db config error")
	}
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte(JsonNull), nil
	}
	return json.Marshal(ni.Int64)
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte(JsonNull), nil
	}
	return json.Marshal(nf.Float64)
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte(JsonNull), nil
	}
	return json.Marshal(ns.String)
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte(JsonNull), nil
	}
	return json.Marshal(nb.Bool)
}

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	mysql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte(JsonNull), nil
	}
	return json.Marshal(dateTimeUtil.Format(nt.Time))
}
