GIT_TOP_LEVEL= $(shell git rev-parse --show-toplevel)
cd ${GIT_TOP_LEVEL}
go run cmd/news/main.go