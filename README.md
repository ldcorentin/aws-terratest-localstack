# aws-terratest-localstack
Test your terraform modules and code using terratest and localstack

LocalStack is a cloud service emulator that runs in a single container on your laptop or in your CI environment. With LocalStack, you can run your AWS applications or Lambdas entirely on your local machine without connecting to a remote cloud provider

## Prerequisites

### Install golang
```
1) curl -o golang.pkg https://dl.google.com/go/go1.16.4.darwin-amd64.pkg
2) sudo open golang.pkg
3) export GOROOT=/usr/local/go
4) export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
5) go version
```
### Install localstack

If you want to use the local stack cli run : `pip3 install localstack`

If you prefer to use directly docker or docker-compose, install them if you don't have it :D 

### Testing

#### Launch localstack
1) Using localstack cli: `SERVICES=s3,sts localstack start`
2) Using localstack cli in detached mode (running in background): `SERVICES=s3,sts localstack start -d`
3) Using docker-compose: `cd docker && docker-compose up` or `docker-compose -f docker/docker-compose.yml up`
4) Using docker-compose in detached mode (running in background): `cd docker && docker-compose up -d` or `docker-compose -f docker/docker-compose.yml up -d`

To check that localstack is running properly, please run `curl http://localhost:4566` and you should get an answer like : `{"status": "running"}`


#### Run the test
```
1) cd test/terratest
2) go mod tidy
3) go test -v
```