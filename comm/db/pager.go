package db

import (
	"context"
	"reflect"
)

type order interface {
	GetOrderType() int32
	GetOrderCol() string
}

type LimitOffset interface {
	GetOffset() int32
	GetLimit() int32
}

type LimitPage interface {
	GetPageSize() int32
	GetPageNo() int32
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
