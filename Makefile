.PHONY: mock
mock:
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/web/mock/user_gen.go
	@go mod tidy