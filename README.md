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

