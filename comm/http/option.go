package http

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Option func(fields gin.H)

func StatusOption(status int) Option {
	return func(fields gin.H) {
		fields["code"] = status
	}
}

func TraceOption(requestId interface{}) Option {
	return func(fields gin.H) {
		fields["request_id"] = requestId
	}
}

func DataOption(data interface{}) Option {
	return func(fields gin.H) {
		fields["data"] = data
	}
}

func FlatOption(target interface{}) Option {
	var data gin.H
	jmv, _ := json.Marshal(target)
	json.Unmarshal(jmv, &data)
	return func(fields gin.H) {
		for k, v := range data {
			fields[k] = v
		}
	}
}

func MsgOption(format string, args ...interface{}) Option {
	return func(fields gin.H) {
		fields["msg"] = fmt.Sprintf(format, args...)
	}
}
