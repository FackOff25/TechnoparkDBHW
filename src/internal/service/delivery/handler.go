package delivery

import (
	serviceUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/service/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Delivery struct {
	ServiceUC serviceUsecase.UseCaseInterface
}

func MakeDelivery(e *echo.Echo, serviceUC serviceUsecase.UseCaseInterface) {
	handler := &Delivery{
		ServiceUC: serviceUC,
	}

	e.POST("/api/service/clear", handler.ClearData)
	e.GET("/api/service/status", handler.SelectStatus)
}

func (delivery *Delivery) ClearData(c echo.Context) error {
	err := delivery.ServiceUC.ClearData()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (delivery *Delivery) SelectStatus(c echo.Context) error {
	status, err := delivery.ServiceUC.SelectStatus()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}
