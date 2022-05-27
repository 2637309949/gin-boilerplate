## gin-boilerplate

### Login

```shell
curl -X POST \
  http://192.168.202.128:8080/api/v1/user/login \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: b4ba87e4-e139-0ea5-f8ba-91a60ec37a77' \
  -d '{
	"email": "2637309949@qq.com",
	"password": "admin"
}'
```

### Refresh

```shell
curl -X POST \
  http://192.168.202.128:8080/api/v1/token/refresh \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: b1763bb7-0d1a-381f-5502-3da04bf23f19' \
  -d '{
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQyMzk5ODcsInJlZnJlc2hfdXVpZCI6IkJaT0xDTXI3Z3oiLCJ1c2VyX2lkIjoxfQ.4sY7PkBdw1PTap-zRdIaWHj8uG_6SrY0zbzDFsPJUKE"
}'
```
