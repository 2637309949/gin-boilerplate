package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Option struct {
	Dialect     string
	DNS         string
	Index       gin.HandlerFunc
	NoRoute     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
	Static      struct {
		RelativePath string
		Root         string
	}
	Swagger string
	Metrics string
	Sql     string
}

type OptFunc func(o *Option)

//DataBase...
func DataBase(dialect string, dns string) OptFunc {
	return func(o *Option) {
		o.DNS = dns
		o.Dialect = dialect
	}
}

//DataBase...
func Swagger(dir string) OptFunc {
	return func(o *Option) {
		o.Swagger = dir
	}
}

//Sql...
func Sql(file string) OptFunc {
	return func(o *Option) {
		o.Sql = file
	}
}

//Mode...
func Mode(value string) OptFunc {
	return func(o *Option) {
		gin.SetMode(value)
	}
}

//Validator...
func Validator(validator binding.StructValidator) OptFunc {
	return func(o *Option) {
		binding.Validator = validator
	}
}

//Middleware...
func Middleware(middles ...gin.HandlerFunc) OptFunc {
	return func(o *Option) {
		o.Middlewares = append(o.Middlewares, middles...)
	}
}

//Index...
func Index(i gin.HandlerFunc) OptFunc {
	return func(o *Option) {
		o.Index = i
	}
}

//NoRoute...
func NoRoute(n gin.HandlerFunc) OptFunc {
	return func(o *Option) {
		o.NoRoute = n
	}
}

//Static...
func Static(relativePath string, root string) OptFunc {
	return func(o *Option) {
		o.Static.RelativePath = relativePath
		o.Static.Root = root
	}
}

//Metrics...
func Metrics(m string) OptFunc {
	return func(o *Option) {
		o.Metrics = m
	}
}
