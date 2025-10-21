package controller

import "github.com/labstack/echo/v4"

type Server interface {
	Start(string) error
	DefineRoutes() error
}

type serverImpl struct {
	server *echo.Echo
}

func NewServer() Server {
	return serverImpl{
		server: echo.New(),
	}
}

func (ref serverImpl) Start(address string) error {
	if err := ref.server.Start(address); err != nil {
		return err
	}
	return nil
}

func (ref serverImpl) DefineRoutes() error {
	return nil
}
