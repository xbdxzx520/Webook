// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/account/grpc"
	"gitee.com/geekbang/basic-go/webook/account/ioc"
	"gitee.com/geekbang/basic-go/webook/account/repository"
	"gitee.com/geekbang/basic-go/webook/account/repository/dao"
	"gitee.com/geekbang/basic-go/webook/account/service"
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
)

// Injectors from wire.go:

func Init() *wego.App {
	db := ioc.InitDB()
	accountDAO := dao.NewCreditGORMDAO(db)
	accountRepository := repository.NewAccountRepository(accountDAO)
	accountService := service.NewAccountService(accountRepository)
	accountServiceServer := grpc.NewAccountServiceServer(accountService)
	client := ioc.InitEtcdClient()
	loggerV1 := ioc.InitLogger()
	server := ioc.InitGRPCxServer(accountServiceServer, client, loggerV1)
	app := &wego.App{
		GRPCServer: server,
	}
	return app
}
