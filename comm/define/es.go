package define

import (
	"fmt"

	"github.com/olivere/elastic"
)

type Query struct {
	i interface{}
}

func (q *Query) Source() (interface{}, error) {
	return q.i, nil
}

func NewSource(i interface{}) elastic.Query {
	return &Query{i}
}

func NewSourcef(i interface{}, a ...interface{}) elastic.Query {
	switch j := i.(type) {
	case string:
		i = fmt.Sprintf(j, a...)
	}
	return &Query{i}
}

func NewScript(i interface{}) elastic.Query {
	return NewSource(map[string]interface{}{
		"script": map[string]interface{}{
			"script": i,
		},
	})
}

func NewScriptf(i interface{}, a ...interface{}) elastic.Query {
	switch j := i.(type) {
	case string:
		i = fmt.Sprintf(j, a...)
	}
	return NewSource(map[string]interface{}{
		"script": map[string]interface{}{
			"script": i,
		},
	})
}
