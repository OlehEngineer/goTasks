package model

import (
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type handler struct{}

func RegisterRouters(e *echo.Echo) {
	h := handler{}
	e.GET("/api/v1/users", h.getUsersPagination)
	e.GET("/api/v1/users/:id", h.getUser)
	e.DELETE("/api/v1/users/:id", h.deleteUser)
	e.POST("/api/v1/users", h.postUser)
	e.PUT("/api/v1/users/:id", h.updateUser)

}

// get information about all users. Pagination implemented.
func (h *handler) getUsersPagination(c echo.Context) error {
	user := []ApiResponse{}

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
	log.Infof("page - %v. Limit - %v", page, limit)
	// get connection to database
	conn := c.Get("conn").(*sqlx.DB)

	totalPageQty, err := GetPageQty(conn, limit)
	if err != nil {
		log.Errorf("getUsersPagination. cannot get page q-ty. error - %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if page > totalPageQty || page <= 0 {
		log.Errorf("getUsersPagination. incorrect page number.")
		return echo.NewHTTPError(http.StatusInternalServerError, "page number out of the available page range")
	}

	user, err = GetUsersPage(conn, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// get information about certain user. Basic Authentication implemented
func (h *handler) getUser(c echo.Context) error {
	user := ApiResponse{}
	id := c.Param("id") //get user's ID from endpoint

	//convert user's ID to uint16
	userID, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// get connection to database
	conn := c.Get("conn").(*sqlx.DB)

	// get user's credentials for Basic Auth verification
	nickname, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid credentials")
	}

	//perform Authentication
	authorized, err := Authentication(conn, nickname, password, uint16(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if authorized != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// return user data after successful verification
	user, errGet := GetUser(conn, uint16(userID))
	if errGet != nil {
		log.Errorf("cannot get the user from the database. Error - %v", errGet)

		return c.JSON(http.StatusNotFound, map[string]string{
			"error":   errGet.Error(),
			"message": "there is no user with this ID in the database",
		})
	}
	return c.JSON(http.StatusOK, user)
}

// create new user. Basic Authentication implemented
func (h *handler) postUser(c echo.Context) error {
	newUser := &UserSignUp{}

	//parsing provided JSON data
	if err := c.Bind(newUser); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	//crypt password
	cryptedPassword, errPass := PasswordHashing(newUser.Password)
	if errPass != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPass.Error())
	}

	//check if provided credentials not nil
	if newUser.NickName == "" || newUser.Name == "" || newUser.LastName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "not all required files available")
	}
	conn := c.Get("conn").(*sqlx.DB)

	//create user inside the database
	userInfo, postErr := PostUser(conn, newUser.NickName, newUser.Name, newUser.LastName, cryptedPassword)
	if postErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, postErr.Error())
	}

	return c.JSON(http.StatusCreated, userInfo)
}

// delete user. Basic Authentication implemented
func (h *handler) deleteUser(c echo.Context) error {

	id := c.Param("id")
	//convert user's ID to uint16
	userID, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		log.Errorf("deleteUser, userID convert to UINT16 error - %v", err)
		return err
	}

	// get connection to database
	conn := c.Get("conn").(*sqlx.DB)

	// get user's credentials for Basic Auth verification
	nickname, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid credentials")
	}

	//perform Authentication
	authorized, err := Authentication(conn, nickname, password, uint16(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if authorized != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	//delete user from database
	errDel := DeleteUser(conn, uint16(userID))
	if errDel != nil {
		log.Errorf("user delete problem. Error - %v", errDel)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   errDel.Error(),
			"message": "cannot delete the user. Internal server error.",
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// update user's information. Basic Authentication implemented
func (h *handler) updateUser(c echo.Context) error {
	updatedUser := &UpdateUser{}
	id := c.Param("id")
	//convert user's ID to uint16
	userID, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		log.Errorf("updateUser, userID convert to UINT16 error - %v", err)
		return err
	}

	// get connection to database
	conn := c.Get("conn").(*sqlx.DB)

	// get user's credentials for Basic Auth verification
	nickname, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid credentials")
	}

	//perform Authentication
	authorized, err := Authentication(conn, nickname, password, uint16(userID))
	if err != nil {
		log.Errorf("Authentication problem. Error - %s", err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if authorized != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	//parsing provided JSON data
	if err := c.Bind(updatedUser); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}
	// update user data
	upUserResponse, err := PutUser(conn, *updatedUser)
	if err != nil {
		log.Errorf("user update error - %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, upUserResponse)
}
