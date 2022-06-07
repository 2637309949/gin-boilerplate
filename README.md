## gin-boilerplate

```shell
[double@double gin-boilerplate]$ go run main.go 
2022-06-07 11:12:21  file=gen/gen.go:76 level=info Generate swagger docs....
2022-06-07 11:12:21  file=web/web.go:47 level=info Listen and serve 380029/main on 0.0.0.0:8080
2022-06-07 11:12:21  file=web/web.go:48 level=info Exec `kill -1 380029` to graceful upgrade
```
### Feature

- TraceId/RequestId for all access
- Standard restful response
- Generic store interface
- Generic pubsub interface
- Automatically generate Swagger documents
- Metrics interface
- Zero downtime restarts or upgrades
- Time-consuming calculation
- Logger with traceId and tracePath

#### Login

```shell
curl -X POST \
http://192.168.202.128:8080/api/v1/user/login \
-H 'content-type: application/json' \
-d '{
  "email": "2637309949@qq.com",
  "password": "admin"
}'
```

#### Refresh

```shell
curl -X POST \
http://192.168.202.128:8080/api/v1/token/refresh \
-H 'content-type: application/json' \
-d '{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}'
```
