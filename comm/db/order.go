package db

const (
	ORDER_NONE = iota
	ORDER_ASC
	ORDER_DESC
)

type order interface {
	GetOrderType() int32
	GetOrderCol() string
}
