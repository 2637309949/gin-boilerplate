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
