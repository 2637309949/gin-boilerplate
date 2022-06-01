package types

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

//OptionsetForm ...
type OptionsetForm struct {
	OptionsetValidator
	Value string `form:"value" json:"value" binding:"required"`
	Code  string `form:"code" json:"code" binding:"required"`
	Name  string `form:"name" json:"name"`
}

//OptionsetFilter...
type OptionsetFilter struct {
	OptionsetValidator
	PageNo    int32  `form:"page_no" json:"page_no"`
	PageSize  int64  `form:"page_size" json:"page_size"`
	OrderType int32  `form:"order_type" json:"order_type"`
	OrderCol  string `form:"order_col" json:"order_col"`
	Name      string `form:"name" json:"name" binding:"required"`
}

func (m *OptionsetFilter) GetPageNo() int32 {
	if m != nil {
		return m.PageNo
	}
	return 0
}

func (m *OptionsetFilter) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *OptionsetFilter) GetOrderType() int32 {
	if m != nil {
		return m.OrderType
	}
	return 0
}

func (m *OptionsetFilter) GetOrderCol() string {
	if m != nil {
		return m.OrderCol
	}
	return ""
}

//OptionsetValidator ...
type OptionsetValidator struct{}

//Name ...
func (f OptionsetValidator) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter optionset name"
		}
		return errMsg[0]
	default:
		return "Something went wrong, please try again later"
	}
}

//Signin ...
func (f OptionsetValidator) Insert(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return f.Name(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//Signin ...
func (f OptionsetValidator) Filter(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return f.Name(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}
