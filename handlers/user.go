package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Shubhouy1/asset-management/database"
	"github.com/Shubhouy1/asset-management/database/dbhelpers"
	"github.com/Shubhouy1/asset-management/middleware"
	"github.com/Shubhouy1/asset-management/models"
	"github.com/Shubhouy1/asset-management/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

var validate = validator.New()

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body models.UserRequest
	var userSessionId string
	var userId string
	var userRole string
	if err := utils.ParseBody(r, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to parse body", err)
		return
	}
	validateErr := validate.Struct(&body)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to validate body", validateErr)
		return
	}

	exists, err := dbhelpers.IsUserExist(body.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to check user", err)
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, "fail to create user", errors.New("user already exists"))
		return
	}
	hashPassword, hashErr := utils.HashPassword(body.Password)
	if hashErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to hash password", hashErr)
		return
	}
	joiningDate, err := time.Parse("2006-01-02", body.JoiningDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid joining_date format (YYYY-MM-DD)", nil)
		return
	}
	userRole = body.Role
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userId, err = dbhelpers.CreateUser(tx, body.Name, body.Email, body.Role, body.Type, body.PhoneNumber, hashPassword, joiningDate)
		if err != nil {
			return err
		}
		userSessionId, err = dbhelpers.CreateUserSession(tx, userId)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to create user", txErr)
		return
	}
	token, err := utils.GenerateJWT(userId, userSessionId, userRole)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to generate token", err)
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"Message":     "user created",
		"accessToken": token,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	var userSessionId string
	var userId string
	var userRole string
	if err := utils.ParseBody(r, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to parse body", err)
		return
	}
	validErr := validate.Struct(&body)
	if validErr != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to validate body", validErr)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var err error
		userId, userRole, err = dbhelpers.GetUserByEmail(tx, body.Email, body.Password)
		if err != nil {
			return err
		}
		userSessionId, err = dbhelpers.CreateUserSession(tx, userId)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusUnauthorized, "invalid email or password", txErr)
		return
	}
	token, err := utils.GenerateJWT(userId, userSessionId, userRole)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to generate token", err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Message":     "user logged in",
		"accessToken": token,
		"userRole":    userRole,
	})
}
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuthContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	sessionID := auth.SessionID
	err := dbhelpers.ArchivedSession(sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "logout failed", err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Message": "user logged out",
	})

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	name := query.Get("name")
	role := query.Get("role")
	userType := query.Get("type")
	assetStatus := query.Get("status")

	users, err := dbhelpers.GetUserInfo(name, role, userType, assetStatus)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch users", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
	})
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "id")

	if userId == "" {
		utils.RespondError(w, http.StatusBadRequest, "user id required", nil)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelpers.ReturnAllAssets(tx, userId)
		if err != nil {
			return err
		}
		err = dbhelpers.ArchiveUserSession(tx, userId)
		if err != nil {
			return err
		}
		err = dbhelpers.DeleteUser(tx, userId)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete user", txErr)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "user deleted successfully",
	})
}
