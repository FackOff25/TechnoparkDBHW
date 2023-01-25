package delivery

import (
	"net/http"
	"strconv"

	forumUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/usecase"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Delivery struct {
	ForumUC forumUsecase.UseCaseInterface
}

func MakeDelivery(e *echo.Echo, forumUC forumUsecase.UseCaseInterface) {
	handler := &Delivery{
		ForumUC: forumUC,
	}

	e.POST("/api/forum/create", handler.CreateForum)
	e.GET("/api/forum/:slug/details", handler.SelectForum)
	e.GET("/api/forum/:slug/users", handler.SelectForumUsers)
}

func (delivery *Delivery) CreateForum(c echo.Context) error {
	var forum models.Forum
	err := c.Bind(&forum)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	err = delivery.ForumUC.CreateForum(&forum)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		case errors.Is(err, models.ErrConflict):
			return c.JSON(http.StatusConflict, forum)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, forum)
}

func (delivery *Delivery) SelectForum(c echo.Context) error {
	forum, err := delivery.ForumUC.SelectForum(c.Param("slug"))
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, forum)
}

func (delivery *Delivery) SelectForumUsers(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	since := c.QueryParam("since")

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	users, err := delivery.ForumUC.SelectForumUsers(c.Param("slug"), limit, since, desc)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, users)
}
