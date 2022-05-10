package db

import (
	"context"
	"gin-boilerplate/models"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

const (
	ORDER_NONE = iota
	ORDER_ASC
	ORDER_DESC
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
func SetDsn(dialect string, args ...interface{}) {
	var err error
	db, err = gorm.Open(dialect, args...)
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
