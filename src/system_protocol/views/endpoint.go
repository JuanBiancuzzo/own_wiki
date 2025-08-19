package views

import "github.com/labstack/echo/v4"

type Endpoint func(ec echo.Context) error

// Hacer un endpoint para actualizar, insertar y eliminar
