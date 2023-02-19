package app

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"my-rest-api/config"
	"my-rest-api/internal/model"
	"my-rest-api/internal/service"
	"net/http"
	"strconv"
)

const (
	getUsersURL   = "/users"
	getUserURL    = "/users/:id"
	createUserURL = "/users"
	updateUserURL = "/users/:id"
	deleteUserURL = "/users/:id"
)

type Api struct {
	echo        *echo.Echo
	config      *config.Config
	userService service.UserService
	//userHandler handler.UserHandler
}

func NewApi(config *config.Config, userService service.UserService) *Api {
	api := &Api{}
	api.config = config
	api.userService = userService

	api.echo = echo.New()
	api.echo.Pre(middleware.RemoveTrailingSlash())

	api.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "user_agent=${user_agent}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	api.echo.GET(getUsersURL, api.getUsers)
	api.echo.GET(getUserURL, api.getUser)
	api.echo.POST(createUserURL, api.createUser)
	api.echo.PUT(updateUserURL, api.updateUser)
	api.echo.DELETE(deleteUserURL, api.deleteUser)

	return api
}

func (api *Api) Start() error {
	return api.echo.Start(":" + api.config.Listen.Port)
}

func (api *Api) getUsers(c echo.Context) error {

	users, err := api.userService.GetAllUsers()

	if err != nil {
		return err
	}

	if users == nil {
		users = []*model.User{}
	}
	return c.JSON(http.StatusOK, users)
}

func (api *Api) getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}

	user, err := api.userService.GetUserById(id)

	if err != nil {
		return err
		//return echo.NewHTTPError(http.StatusNotAcceptable, "Error: " + err.Error())
	}

	if user == nil {
		errorString := "User with `id` = " + c.Param("id") + " not found"
		return echo.NewHTTPError(http.StatusNotFound, errorString)
	}
	return c.JSON(http.StatusOK, user)
}

func (api *Api) createUser(c echo.Context) error {
	user := &model.User{}

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}

	res, err := api.userService.AddUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]int{"id": res.Id})

}

func (api *Api) updateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}

	user := &model.User{}
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	if err := user.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}

	user.Id = id

	if err = api.userService.UpdateUser(user); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			errorString := "User with `id` = " + c.Param("id") + " not found"
			return echo.NewHTTPError(http.StatusNotFound, errorString)
		}

	}
	return c.JSON(http.StatusOK, "User is updated")
}

func (api *Api) deleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}

	if err = api.userService.DeleteUser(id); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			errorString := "User with `id` = " + c.Param("id") + " not found"
			return echo.NewHTTPError(http.StatusNotFound, errorString)
		}

	}
	return c.JSON(http.StatusOK, "User is deleted")

}
