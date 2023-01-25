package delivery

import (
	"net/http"

	userUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/user/usecase"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Delivery struct {
	UserUC userUsecase.UseCaseInterface
}

func MakeDelivery(e *echo.Echo, userUC userUsecase.UseCaseInterface) {
	handler := &Delivery{
		UserUC: userUC,
	}

	e.POST("/api/user/:nickname/create", handler.CreateUser)
	e.GET("/api/user/:nickname/profile", handler.SelectUser)
	e.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

func (delivery *Delivery) CreateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user.NickName = c.Param("nickname")

	conflictUsers, err := delivery.UserUC.CreateUser(&user)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrConflict):
			return c.JSON(http.StatusConflict, conflictUsers)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func (delivery *Delivery) SelectUser(c echo.Context) error {
	user, err := delivery.UserUC.SelectUser(c.Param("nickname"))
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (delivery *Delivery) UpdateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user.NickName = c.Param("nickname")

	err = delivery.UserUC.UpdateUser(&user)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		case errors.Is(err, models.ErrConflict):
			return echo.NewHTTPError(http.StatusConflict, models.ErrConflict.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, user)
}
