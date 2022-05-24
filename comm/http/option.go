package http

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Option func(v gin.H)

func StatusOption(status int) Option {
	return func(v gin.H) {
		v["code"] = status
	}
}

func TraceOption(requestId interface{}) Option {
	return func(v gin.H) {
		v["request_id"] = requestId
	}
}

func DataOption(data interface{}) Option {
	return func(v gin.H) {
		v["data"] = data
	}
}

func FlatOption(target interface{}) Option {
	var data gin.H
	jmv, _ := json.Marshal(target)
	json.Unmarshal(jmv, &data)
	return func(v gin.H) {
		for k, v1 := range data {
			v[k] = v1
		}
	}
}

func MsgOption(format string, args ...interface{}) Option {
	return func(v gin.H) {
		v["msg"] = fmt.Sprintf(format, args...)
	}
}
