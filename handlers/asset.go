package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Shubhouy1/asset-management/database"
	"github.com/Shubhouy1/asset-management/database/dbhelpers"
	"github.com/Shubhouy1/asset-management/middleware"
	"github.com/Shubhouy1/asset-management/models"
	"github.com/Shubhouy1/asset-management/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {

	var body models.CreateAssetRequest
	var assetID string

	if err := utils.ParseBody(r, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body", nil)
		return
	}

	if err := validate.Struct(&body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "validation failed", nil)
		return
	}

	warrantyStart, err := time.Parse("2006-01-02", body.WarrantyStart)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid warrantyStart", nil)
		return
	}

	warrantyEnd, err := time.Parse("2006-01-02", body.WarrantyEnd)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid warrantyEnd", nil)
		return
	}

	if warrantyEnd.Before(warrantyStart) {
		utils.RespondError(w, http.StatusBadRequest, "invalid warranty range", nil)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {

		var err error
		assetID, err = dbhelpers.CreateAsset(tx, body.Brand, body.Model, body.SerialNo, body.Type, body.Owner, warrantyStart, warrantyEnd)
		if err != nil {
			return err
		}
		switch body.Type {

		case "laptop":
			if body.Laptop == nil {
				return fmt.Errorf("laptop details required")
			}
			return dbhelpers.InsertLaptop(tx, assetID, body.Laptop)

		case "mouse":
			if body.Mouse == nil {
				return fmt.Errorf("mouse details required")
			}
			return dbhelpers.InsertMouse(tx, assetID, body.Mouse)
		case "keyboard":
			if body.Keyboard == nil {
				return fmt.Errorf("keyboard details required")
			}
			return dbhelpers.InsertKeyboard(tx, assetID, body.Keyboard)
		case "mobile":
			if body.Mobile == nil {
				return fmt.Errorf("mobile details required")
			}
			return dbhelpers.InsertMobile(tx, assetID, body.Mobile)

		default:
			return fmt.Errorf("unsupported asset type")
		}
	})

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create asset", nil)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"assetId": assetID,
	})
}
func AssignAsset(w http.ResponseWriter, r *http.Request) {
	assetId := chi.URLParam(r, "id")
	var body models.AssignRequest
	err := utils.ParseBody(r, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body", nil)
		return
	}
	validateErr := validate.Struct(&body)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to validate body", validateErr)
		return
	}
	auth, ok := middleware.GetAuthContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userId := auth.UserID
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var err error
		err = dbhelpers.AssignAsset(tx, assetId, body.AssignedTo, userId)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to assign asset", txErr)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Message": "asset assigned successfully",
	})
}
func TotalAssets(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuthContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userId := auth.UserID
	totalAssets, err := dbhelpers.FindTotalAssetById(userId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to find total assets", err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"totalAssets": totalAssets,
	})
}
func SentToService(w http.ResponseWriter, r *http.Request) {
	assetId := chi.URLParam(r, "id")
	if assetId == "" {
		utils.RespondError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}
	var body models.SentServiceRequest
	err := utils.ParseBody(r, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body", nil)
		return
	}
	err = validate.Struct(&body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "validation failed", nil)
		return
	}

	serviceStart, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid service date", nil)
		return
	}
	serviceEnd, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid service end date", nil)
		return
	}
	if serviceEnd.Before(serviceStart) {
		utils.RespondError(w, http.StatusBadRequest, "end date must be after start date", nil)
		return
	}
	err = dbhelpers.SentToService(assetId, serviceStart, serviceEnd)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "fail to sent for service", err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "asset sent for service successfully",
	})

}
func ShowAssets(w http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")
	statusStr := r.URL.Query().Get("status")
	ownerStr := r.URL.Query().Get("owner")
	brandStr := r.URL.Query().Get("brand")
	modelStr := r.URL.Query().Get("model")
	serialNumberStr := r.URL.Query().Get("serialNumber")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 5

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p <= 0 {
			utils.RespondError(w, http.StatusBadRequest, "invalid page", nil)
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			utils.RespondError(w, http.StatusBadRequest, "invalid limit", nil)
			return
		}
		limit = l
	}
	offset := (page - 1) * limit

	assetsData, err := dbhelpers.ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr, limit, offset)
	if err != nil {

		utils.RespondError(w, http.StatusInternalServerError, "fail to fetch assets", err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assetsData,
	})
}
func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	assetId := chi.URLParam(r, "id")
	if assetId == "" {
		utils.RespondError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}
	var body models.UpdateAssetRequest
	err := utils.ParseBody(r, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body", nil)
		return
	}
	validateErr := validate.Struct(&body)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to validate body", validateErr)
		return
	}
	warrantyStart, err := time.Parse("2006-01-02", body.WarrantyStart)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid warrantyStart", nil)
		return
	}

	warrantyEnd, err := time.Parse("2006-01-02", body.WarrantyEnd)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid warrantyEnd", nil)
		return
	}

	if warrantyEnd.Before(warrantyStart) {
		utils.RespondError(w, http.StatusBadRequest, "invalid warranty range", nil)
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelpers.UpdateAsset(tx, assetId, body.Brand, body.Model, body.SerialNo, body.Type, body.Owner, warrantyStart, warrantyEnd)
		if err != nil {
			return err
		}
		switch body.Type {

		case "laptop":
			if body.Laptop == nil {
				return fmt.Errorf("laptop details required")
			}
			return dbhelpers.UpdateLaptop(tx, assetId, body.Laptop)

		case "mouse":
			if body.Mouse == nil {
				return fmt.Errorf("mouse details required")
			}
			return dbhelpers.UpdateMouse(tx, assetId, body.Mouse)
		case "keyboard":
			if body.Keyboard == nil {
				return fmt.Errorf("keyboard details required")
			}
			return dbhelpers.UpdateKeyboard(tx, assetId, body.Keyboard)
		case "mobile":
			if body.Mobile == nil {
				return fmt.Errorf("mobile details required")
			}
			return dbhelpers.UpdateMobile(tx, assetId, body.Mobile)

		default:
			return fmt.Errorf("unsupported asset type")
		}
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusBadRequest, "fail to update asset", txErr)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "asset updated",
	})
}
