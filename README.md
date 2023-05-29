# Sample Echo Golang API

# Build
`make binary`

# Docker commands for go-echoapp

## Docker build 
`docker build --tag goechoapp .`
## Docker tag image
`docker image tag goechoapp:latest goechoapp:v1.0`

`docker image tag goechoapp:latest akoserwal/goechoapp:v1.0`

## docker run
`docker run --publish 3000:3000 go-echoapp`

## docker run detached mode
`docker run -d -p 3000:3000 go-echoapp`