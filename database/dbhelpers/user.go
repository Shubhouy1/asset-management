package dbhelpers

import (
	"time"

	"github.com/Shubhouy1/asset-management/database"
	"github.com/Shubhouy1/asset-management/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func GetUserIDFromSession(sessionId string) (userId string, err error) {
	query := `Select user_id 
             from user_session where id=$1
             And archived_at is null`

	err = database.Asset.Get(&userId, query, sessionId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
func GetUserByEmail(tx *sqlx.Tx, email, password string) (string, string, error) {
	query := `Select id, password_hash,role
             from users 
             where TRIM(LOWER(email)) = LOWER($1)
             And archived_at is null`
	var result models.UserExist
	err := tx.Get(&result, query, email)
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(result.PasswordHash),
		[]byte(password)); err != nil {
		return "", "", err
	}
	return result.ID, result.Role, nil
}
func IsUserExist(email string) (bool, error) {
	query := `Select Count(*)>0
             from users 
             where trim(LOWER(email)) = trim(lower($1))
             And archived_at is null`
	var exist bool
	err := database.Asset.Get(&exist, query, email)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func CreateUser(tx *sqlx.Tx, name, email, role, userType, phoneNo, passwordHash string, joiningDate time.Time) (string, error) {
	query := `INSERT INTO users (name, email,role,type,phone_no,password_hash,joining_date)
               values ($1, trim(lower($2)), $3, $4, $5, $6, $7) RETURNING id`
	var userId string
	err := tx.Get(&userId, query, name, email, role, userType, phoneNo, passwordHash, joiningDate)
	if err != nil {
		return "", err
	}
	return userId, nil
}
func CreateUserSession(tx *sqlx.Tx, userId string) (string, error) {
	query := `INSERT INTO user_session (user_id)
              values ($1)RETURNING id`
	var userSessionId string
	err := tx.Get(&userSessionId, query, userId)
	if err != nil {
		return "", err
	}
	return userSessionId, nil
}

func ArchivedSession(sessionID string) error {
	query := `Update user_session 
             set archived_at=Now() where id=$1 and archived_at is null`
	_, err := database.Asset.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}
func CreateAsset(tx *sqlx.Tx, brand, model, serialNo, assetType, owner string, warrantyStart, warrantyEnd time.Time) (string, error) {
	query := `Insert into assets (brand, model, serial_no,type, owner,warranty_start, warranty_end)
	          Values($1, $2, $3, $4, $5, $6, $7) returning id`
	var assetId string
	err := tx.Get(&assetId, query, brand, model, serialNo, assetType, owner, warrantyStart, warrantyEnd)
	if err != nil {
		return "", err
	}
	return assetId, nil

}
func AssignAsset(tx *sqlx.Tx, assetId, assignedTo, assignedBy string) error {
	query := `Update assets set
             status ='assigned',
             assigned_to = $2 ,
             assigned_by_id = $3,
             assigned_on = now(),
             updated_at = now()
             where id = $1
             and archived_at is null
             and status = 'available'`
	_, err := tx.Exec(query, assetId, assignedTo, assignedBy)
	if err != nil {
		return err
	}
	return nil
}
func InsertLaptop(tx *sqlx.Tx, assetId string, laptop *models.LaptopInput) error {
	query := `Insert into laptop(asset_id,processor, ram,storage,os, charger, password)
             values ($1, $2, $3, $4, $5, $6,$7)`

	_, err := tx.Exec(query, assetId, laptop.Processor, laptop.RAM, laptop.Storage, laptop.OS, laptop.Charger, laptop.Password)
	if err != nil {
		return err
	}
	return nil
}
func InsertMouse(tx *sqlx.Tx, assetId string, mouse *models.MouseInput) error {
	query := `Insert into mouse(asset_id, dpi, connectivity) 
             values ($1, $2, $3)`
	_, err := tx.Exec(query, assetId, mouse.Dpi, mouse.Connectivity)
	if err != nil {
		return err
	}
	return nil
}
func InsertKeyboard(tx *sqlx.Tx, assetId string, keyboard *models.KeyboardInput) error {
	query := ` Insert into keyboard(asset_id, layout, connectivity)
  values($1, $2, $3)`
	_, err := tx.Exec(query, assetId, keyboard.Layout, keyboard.Connectivity)
	if err != nil {
		return err
	}
	return nil
}

func InsertMobile(tx *sqlx.Tx, assetId string, mobile *models.MobileInput) error {
	query := `Insert into mobile(asset_id, os,ram,storage,charger, password)
              values ($1, $2, $3, $4, $5, $6)`
	_, err := tx.Exec(query, assetId, mobile.Os, mobile.Ram, mobile.Storage, mobile.Charger, mobile.Password)
	if err != nil {
		return err
	}
	return nil
}
func FindTotalAssetById(userID string) (models.DashboardUserData, error) {
	query := `Select count(*) as active_asset from assets 
                where assigned_to=$1 and status='assigned'
                AND archived_at is null`
	var dummy models.DashboardUserData
	var summary models.DashboardUserSummary
	err := database.Asset.Get(&summary, query, userID)
	if err != nil {
		return dummy, err
	}
	assetInfo := make([]models.Asset, 0)
	query1 := `Select brand,model,serial_no,type,status,owner,created_at
               from assets 
               where assigned_to=$1 and archived_at is null`
	err = database.Asset.Select(&assetInfo, query1, userID)
	if err != nil {
		return dummy, err
	}
	return models.DashboardUserData{
		Summary: summary,
		Assets:  assetInfo,
	}, nil
}
func SentToService(assetId string, serviceStart, serviceEnd time.Time) error {
	query := `update assets set status='in_service',service_start=$2,service_end=$3,updated_at=now()
              where id=$1 and archived_at is null and status ='available'`
	_, err := database.Asset.Exec(query, assetId, serviceStart, serviceEnd)
	if err != nil {
		return err
	}
	return nil

}
func ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr string, limit, offset int) (models.DashboardData, error) {
	SQL := `SELECT brand, model, type, serial_no, status, owner, created_at
			FROM assets
			WHERE archived_at IS NULL
			AND ($1 = '' OR brand ILIKE '%'||$1||'%')
			AND ($2 = '' OR model ILIKE '%'||$2||'%')
			AND ($3 = '' OR serial_no ILIKE '%'||$3||'%')
			AND ($4 = '' OR type::text ILIKE '%'||$4||'%')
			AND ($5 = '' OR status::text ILIKE '%'||$5||'%')
			AND ($6 = '' OR owner::text ILIKE '%'||$6||'%')
			ORDER BY created_at DESC
			LIMIT $7 OFFSET $8
          `

	assets := make([]models.Asset, 0)
	var summary models.DashboardSummary

	Sql := `SELECT
          COUNT(*) AS total,
          COUNT(*) FILTER (WHERE status = 'available') AS available,
          COUNT(*) FILTER (WHERE status = 'assigned') AS assigned,
          COUNT(*) FILTER (WHERE status = 'for_repair') AS waiting_for_repair,
          COUNT(*) FILTER (WHERE status = 'in_service') AS in_service,
          COUNT(*) FILTER (WHERE status = 'damaged') AS damaged
       FROM assets
       WHERE archived_at IS NULL`

	var res models.DashboardData
	DashboardErr := database.Asset.Get(&summary, Sql)
	if DashboardErr != nil {
		return res, DashboardErr
	}

	err := database.Asset.Select(&assets, SQL, brandStr, modelStr, serialNumberStr, typeStr, statusStr, ownerStr, limit, offset)
	if err != nil {
		return res, err
	}
	return models.DashboardData{
		Summary: summary,
		Assets:  assets,
	}, nil
}
func UpdateAsset(tx *sqlx.Tx, assetID, brand, model, serialNo, assetType, owner string, warrantyStart, warrantyEnd time.Time) error {
	query := `UPDATE assets
            set brand = $2, model = $3, serial_no = $4, type=$5,owner=$6,warranty_start = $7,warranty_end=$8, updated_at =now()
            where id= $1 and archived_at is null `
	_, err := tx.Exec(query, assetID, brand, model, serialNo, assetType, owner, warrantyStart, warrantyEnd)
	if err != nil {
		return err
	}
	return nil

}
func UpdateLaptop(tx *sqlx.Tx, assetID string, laptop *models.LaptopInput) error {
	query := `
	UPDATE laptop
	SET
	    processor = $2,
	    ram = $3,
	    storage = $4,
	    os = $5,
	    charger = $6,
	    password = $7
	WHERE asset_id = $1
	`

	_, err := tx.Exec(query,
		assetID,
		laptop.Processor,
		laptop.RAM,
		laptop.Storage,
		laptop.OS,
		laptop.Charger,
		laptop.Password,
	)

	return err
}
func UpdateMouse(tx *sqlx.Tx, assetID string, mouse *models.MouseInput) error {
	query := `
	UPDATE mouse
	SET
	    dpi = $2,
	    connectivity = $3
	WHERE asset_id = $1
	`

	_, err := tx.Exec(query, assetID, mouse.Dpi, mouse.Connectivity)
	return err
}
func UpdateKeyboard(tx *sqlx.Tx, assetID string, keyboard *models.KeyboardInput) error {
	query := `
	UPDATE keyboard
	SET
	    layout = $2,
	    connectivity = $3
	WHERE asset_id = $1
	`

	_, err := tx.Exec(query, assetID, keyboard.Layout, keyboard.Connectivity)
	return err
}
func UpdateMobile(tx *sqlx.Tx, assetID string, mobile *models.MobileInput) error {
	query := `
	UPDATE mobile
	SET
	    os = $2,
	    ram = $3,
	    storage = $4,
	    charger = $5,
	    password = $6
	WHERE asset_id = $1
	`

	_, err := tx.Exec(
		query,
		assetID,
		mobile.Os,
		mobile.Ram,
		mobile.Storage,
		mobile.Charger,
		mobile.Password,
	)

	return err
}
func GetAssetInfo(userID, assetStatus string) ([]models.AssetInfoRequest, error) {

	query := `
		SELECT id, brand, model, status, type
		FROM assets
		WHERE assigned_to = $1
		AND archived_at IS NULL
		AND ($2 = '' OR status::TEXT = $2)
	`

	assetDetails := make([]models.AssetInfoRequest, 0)

	err := database.Asset.Select(&assetDetails, query, userID, assetStatus)
	return assetDetails, err
}

func GetUserInfo(name, role, userType, assetStatus string) ([]models.UserInfoRequest, error) {

	query := `
		SELECT id, name, email, phone_no, role, type, created_at
		FROM users
		WHERE archived_at IS NULL
		AND ($1 = '' OR name ILIKE '%' || $1 || '%')
		AND ($2 = '' OR role::TEXT = $2)
		AND ($3 = '' OR type::TEXT = $3)
	`

	users := make([]models.UserInfoRequest, 0)

	err := database.Asset.Select(&users, query, name, role, userType)
	if err != nil {
		return users, err
	}

	filteredUsers := make([]models.UserInfoRequest, 0)

	for _, user := range users {

		assetDetails, err := GetAssetInfo(user.ID, assetStatus)
		if err != nil {
			return users, err
		}

		if assetStatus != "available" && len(assetDetails) == 0 {
			continue
		}

		user.AssetDetails = assetDetails
		filteredUsers = append(filteredUsers, user)
	}

	return filteredUsers, nil
}
func DeleteUser(tx *sqlx.Tx, userId string) error {
	query := `Update users set archived_at = now()
             where id = $1 and archived_at is null`
	_, err := tx.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil

}
func CountActiveAssets(tx *sqlx.Tx, userId string) (int, error) {
	query := `Select count(*)
              from assets 
              where assigned_to = $1
              and archived_at IS NULL`
	var count int
	err := tx.Get(&count, query, userId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func ArchiveUserSession(tx *sqlx.Tx, userId string) error {
	query := `update user_session
              set archived_at = now()
              where user_id = $1
              and archived_at is null`
	_, err := tx.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}
func ReturnAllAssets(tx *sqlx.Tx, userId string) error {

	query := `
		UPDATE assets
		SET assigned_to = NULL,
		    assigned_by_id = NULL,
		    assigned_on = NULL,
		    status = 'available',
		    returned_on = now(),
		    updated_at = now()
		WHERE assigned_to = $1
		AND archived_at IS NULL
	`

	_, err := tx.Exec(query, userId)
	return err
}
