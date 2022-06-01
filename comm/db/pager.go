package db

import (
	"context"
	"reflect"
)

type LimitOffset interface {
	GetOffset() int
	GetLimit() int
}

type LimitPage interface {
	GetPageSize() int
	GetPageNo() int
}

func InitPage(ctx context.Context, itf LimitPage) {
	pageSize := itf.GetPageSize()
	pageNo := itf.GetPageNo()

	if pageNo == 0 {
		pageNo = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	fieldPageNo := reflect.ValueOf(itf).Elem().FieldByName("PageNo")
	fieldPageSize := reflect.ValueOf(itf).Elem().FieldByName("PageSize")
	fieldPageNo.SetInt(int64(pageNo))
	fieldPageSize.SetInt(int64(pageSize))
}
