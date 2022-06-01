package db

import (
	"context"
	"gin-boilerplate/models"
	"io/ioutil"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func SetLimit(ctx context.Context, db *gorm.DB, limiter interface{}) *gorm.DB {
	if l, ok1 := limiter.(LimitOffset); ok1 {
		if l.GetOffset() > 0 {
			db = db.Offset(l.GetOffset())
		}
		if l.GetLimit() > 0 {
			db = db.Limit(l.GetLimit())
		}
	} else if l, ok2 := limiter.(LimitPage); ok2 {
		InitPage(ctx, l)
		db = db.Limit(l.GetPageSize())
		db = db.Offset(l.GetPageSize() * (l.GetPageNo() - 1))
	}
	return db
}

func SetOrder(ctx context.Context, db *gorm.DB, o order, tb ...string) *gorm.DB {
	strOrder := o.GetOrderCol()
	if len(strOrder) > 0 {
		if len(tb) > 0 {
			strOrder = tb[0] + "." + strOrder
		}
		switch o.GetOrderType() {
		case ORDER_ASC:
			strOrder += " ASC"
		case ORDER_DESC:
			strOrder += " DESC"
		default:
			strOrder += " ASC"
		}

		db = db.Order(strOrder)
	}

	return db
}

//SetDsn establishes dsn  to database and saves its handler into db *sqlx.DB
func SetDsn(dialector string, dsn string) {
	var err error
	var gd gorm.Dialector
	switch dialector {
	case "sqlite3", "sqlite":
		gd = sqlite.Open(dsn)
	case "mysql":
		gd = mysql.Open(dsn)
	}
	db, err = gorm.Open(gd)
	if err != nil {
		panic(err)
	}
}

//GetDB returns database handler
func GetDB() *gorm.DB {
	return db
}

//AutoMigrate runs gorm auto migration
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Article{})
}

//Exec runs gorm auto migration
func Exec(file string) {
	sqlByte, _ := ioutil.ReadFile(file)
	if len(sqlByte) > 0 {
		GetDB().Exec(string(sqlByte))
	}
}
