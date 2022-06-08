SERVICE_NAME = gin-boilerplate

build:
	GOOS=linux GOARCH=amd64 go build -o $(SERVICE_NAME)

clean:
	rm -rf $(SERVICE_NAME)
