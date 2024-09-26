//go:build wireinject

package startup

import (
	"awesome-bluebook/repository"
	"awesome-bluebook/repository/dao"
	"awesome-bluebook/service"
	"awesome-bluebook/web"
	"github.com/google/wire"
)

var DBProvider = wire.NewSet(InitDB)
var LoggerProvider = wire.NewSet(InitLogger)
var UserHandlerProvider = wire.NewSet(dao.NewGORMUserDAO, repository.NewBasicUserRepository,
	service.NewBasicUserService, web.NewUserHandler)

func NewUserHandler() *web.UserHandler {
	wire.Build(DBProvider, LoggerProvider, UserHandlerProvider)
	return new(web.UserHandler)
}
