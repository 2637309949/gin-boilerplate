package types

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

//ArticleForm ...
type ArticleForm struct {
	ArticleValidator
	Title         string `form:"title" json:"title" binding:"required,min=3,max=100"`
	Introduction  string `form:"introduction" json:"introduction" binding:"required,min=3,max=200"`
	ContentMd     string `form:"content_md" json:"content_md" binding:"required,min=3,max=10000"`
	ContentHtml   string `form:"content_html" json:"content_html" binding:"required,min=3,max=10000"`
	DirectoryHtml string `form:"directory_html" json:"directory_html" binding:"required,min=3,max=10000"`
	Tags          string `form:"tags" json:"tags"`
}

//ArticleFilter...
type ArticleFilter struct {
	ArticleValidator
	PageNo    int32  `form:"page_no" json:"page_no"`
	PageSize  int32  `form:"page_size" json:"page_size"`
	OrderType int32  `form:"order_type" json:"order_type"`
	OrderCol  string `form:"order_col" json:"order_col"`
	Title     string `form:"title" json:"title" binding:"required"`
}

func (m *ArticleFilter) GetPageNo() int32 {
	if m != nil {
		return m.PageNo
	}
	return 0
}

func (m *ArticleFilter) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ArticleFilter) GetOrderType() int32 {
	if m != nil {
		return m.OrderType
	}
	return 0
}

func (m *ArticleFilter) GetOrderCol() string {
	if m != nil {
		return m.OrderCol
	}
	return ""
}

//ArticleValidator ...
type ArticleValidator struct{}

//Name ...
func (f ArticleValidator) Title(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter article title"
		}
		return errMsg[0]
	default:
		return "Something went wrong, please try again later"
	}
}

//Signin ...
func (f ArticleValidator) Insert(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Title" {
				return f.Title(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//Signin ...
func (f ArticleValidator) Filter(err error) string {
	fmt.Println("-----", err)
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Title" {
				return f.Title(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}
