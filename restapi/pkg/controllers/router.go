package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/domian"
	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/service"
)

type httpHandler struct {
	service service.Service
	e       *echo.Echo
}

func New(service service.Service) *httpHandler {
	h := &httpHandler{
		service: service,
		e:       echo.New(),
	}
	return h
}
func (h *httpHandler) Run() {
	h.RegisterRouters()
	h.e.Logger.Fatal(h.e.Start(":8080"))
}

func (h *httpHandler) RegisterRouters() {
	h.e.GET("/api/users", h.getUsersPagination)
	h.e.GET("/api/users/:id", h.getUser)
	h.e.DELETE("/api/users/:id", h.deleteUser)
	h.e.POST("/api/users", h.createUser)
	h.e.PUT("/api/users/:id", h.updateUser)

}

// get information about all users. Pagination implemented.
func (h *httpHandler) getUsersPagination(c echo.Context) error {

	page, err := strconv.Atoi(c.QueryParam("page")) // get page number
	if err != nil {
		log.Errorf("getUsersPagination. cannot get page number. error - %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		log.Errorf("getUsersPagination. cannot get limit value. error - %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalPageQty, err := h.service.GetPageQty(limit)
	if err != nil {
		log.Errorf("getUsersPagination. cannot get page q-ty. error - %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if page > totalPageQty || page <= 0 {
		log.Errorf("getUsersPagination. incorrect page number.")
		return echo.NewHTTPError(http.StatusInternalServerError, "page number out of the available page range")
	}

	user, err := h.service.GetUsersPage(page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// get user by id. Basic Authentication implemented
func (h *httpHandler) getUser(c echo.Context) error {
	id := c.Param("id") //get user's ID from endpoint

	//convert user's ID to uint16
	userID, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// get user's credentials for Basic Auth verification
	nickname, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid credentials")
	}

	//perform Authentication
	authorized, err := h.service.Authentication(nickname, password, uint16(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if authorized != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// return user data after successful verification
	user, errGet := h.service.GetUser(uint16(userID))
	if errGet != nil {
		log.Errorf("cannot get the user from the database. Error - %v", errGet)

		return c.JSON(http.StatusNotFound, map[string]string{
			"error":   errGet.Error(),
			"message": "there is no user with this ID in the database",
		})
	}
	return c.JSON(http.StatusOK, user)
}

// delete user by id. Basic Authentication implemented
func (h *httpHandler) deleteUser(c echo.Context) error {

	id := c.Param("id")
	//convert user's ID to uint16
	userID, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		log.Errorf("deleteUser, userID convert to UINT16 error - %v", err)
		return err
	}

	// get user's credentials for Basic Auth verification
	nickname, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid credentials")
	}

	//perform Authentication
	authorized, err := h.service.Authentication(nickname, password, uint16(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if authorized != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	//delete user from database
	errDel := h.service.DeleteUser(uint16(userID))
	if errDel != nil {
		log.Errorf("user delete problem. Error - %v", errDel)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   errDel.Error(),
			"message": "cannot delete the user. Internal server error.",
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// update user by id. Basic Authentication implemented
func (h *httpHandler) updateUser(c echo.Context) error {

	id := c.Param("id")
	//convert user's ID to uint16
	userID, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		log.Errorf("updateUser, userID convert to UINT16 error - %v", err)
		return err
	}

	// get user's credentials for Basic Auth verification
	nickname, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid credentials")
	}

	//perform Authentication
	authorized, err := h.service.Authentication(nickname, password, uint16(userID))
	if err != nil {
		log.Errorf("Authentication problem. Error - %s", err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if authorized != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	user := &domian.UpdateUser{}
	//parsing provided JSON data
	if err := c.Bind(user); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}
	// update user data
	upUserResponse, err := h.service.UpdateUser(*user)
	if err != nil {
		log.Errorf("user update error - %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, upUserResponse)
}

// create new user. Basic Authentication implemented
func (h *httpHandler) createUser(c echo.Context) error {

	user := &domian.UserSignUp{}
	//parsing provided JSON data
	if err := c.Bind(user); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	//crypt password
	hashedPassword, errPass := h.service.PasswordHashing(user.Password)
	if errPass != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPass.Error())
	}

	//check if provided credentials not nil
	if user.NickName == "" || user.Name == "" || user.LastName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "not all required files available")
	}

	//create user inside the database
	userInfo, postErr := h.service.CreateUser(*user, hashedPassword)
	if postErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, postErr.Error())
	}

	return c.JSON(http.StatusCreated, userInfo)
}
