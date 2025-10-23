package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Handler interface {
	Post(echo.Context) error
	Get(echo.Context) error
	Put(echo.Context) error
	Delete(echo.Context) error
}

type Server interface {
	Start(string) error
}
type serverImpl struct {
	server *echo.Echo
}

func NewServer(instance *echo.Echo, address string) error {
	if err := instance.Start(address); err != nil {
		log.Fatal().Err(err).Msg("fatal error on server, terminating...")
		return err
	}
	return nil
}
