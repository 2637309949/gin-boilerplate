package db

type Option struct {
	Dialect string
	Args    []interface{}
}

type OptFunc func(o *Option)

func Sqlite3(args ...interface{}) OptFunc {
	return func(o *Option) {
		o.Args = args
		o.Dialect = "sqlite3"
	}
}
