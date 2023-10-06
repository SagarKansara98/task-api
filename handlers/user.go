package handlers

import (
	"net/http"
	"task-api/buisness/user"

	"github.com/labstack/echo/v4"
)

var userResource = "User"

type userHandlers struct {
	userService user.User
}

func (u userHandlers) create(c echo.Context) error {
	newUser := user.NewInfo{}
	err := c.Bind(&newUser)
	if err != nil {
		return err
	}

	err = c.Validate(newUser)
	if err != nil {
		return err
	}

	createdUser, err := u.userService.Create(c.Request().Context(), newUser)
	if err != nil {
		if err == user.ErrEmailAlreadyExist {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		return Error(c, err, userResource)
	}

	return c.JSON(http.StatusOK, createdUser)
}

func (u userHandlers) queryByID(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	foundUser, err := u.userService.QueryByID(c.Request().Context(), userID)
	if err != nil {
		return Error(c, err, userResource)
	}

	return c.JSON(http.StatusOK, foundUser)
}

func (u userHandlers) update(c echo.Context) error {
	updateUser := user.Info{}
	err := c.Bind(&updateUser)
	if err != nil {
		return err
	}

	err = c.Validate(updateUser)
	if err != nil {
		return err
	}

	updateUser.ID = getUserID(c)
	if updateUser.ID == 0 {
		return echo.ErrUnauthorized
	}

	err = u.userService.Update(c.Request().Context(), updateUser)
	if err != nil {
		return Error(c, err, userResource)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user updated successfully"})
}

func (u userHandlers) delete(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	err := u.userService.Delete(c.Request().Context(), userID)
	if err != nil {
		return Error(c, err, userResource)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func (u userHandlers) login(c echo.Context) error {
	userInfo := user.Login{}
	err := c.Bind(&userInfo)
	if err != nil {
		return err
	}

	err = c.Validate(&userInfo)
	if err != nil {
		return err
	}

	token, err := u.userService.Login(c.Request().Context(), userInfo)
	if err != nil {
		if err == user.ErrInvalidPassword || err == user.ErrInvalidEmail {
			return badRequest(c, err.Error())
		}

		return Error(c, err, resource)
	}

	return c.JSON(http.StatusOK, map[string]string{"access_token": token})
}
