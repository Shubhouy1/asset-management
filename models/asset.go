package models

import "time"

type LaptopInput struct {
	Processor string `json:"processor" validate:"required"`
	RAM       string `json:"ram" validate:"required"`
	Storage   string `json:"storage" validate:"required"`
	OS        string `json:"os" validate:"required"`
	Charger   string `json:"charger"`
	Password  string `json:"password" validate:"required"`
}
type CreateAssetRequest struct {
	Brand         string `json:"brand" validate:"required"`
	Model         string `json:"model" validate:"required"`
	SerialNo      string `json:"serialNo" validate:"required"`
	Type          string `json:"type" validate:"required"`
	Owner         string `json:"owner" validate:"required"`
	WarrantyStart string `json:"warrantyStart" validate:"required"`
	WarrantyEnd   string `json:"warrantyEnd" validate:"required"`

	Laptop   *LaptopInput   `json:"laptop,omitempty"`
	Mouse    *MouseInput    `json:"mouse,omitempty"`
	Keyboard *KeyboardInput `json:"keyboard,omitempty"`
	Mobile   *MobileInput   `json:"mobile,omitempty"`
}
type MouseInput struct {
	Dpi          int    `json:"dpi" validate:"required"`
	Connectivity string `json:"connectivity" validate:"required"`
}
type KeyboardInput struct {
	Layout       string `json:"layout" validate:"required"`
	Connectivity string `json:"connectivity" validate:"required"`
}
type MobileInput struct {
	Os       string `json:"os" validate:"required"`
	Ram      string `json:"ram" validate:"required"`
	Storage  string `json:"storage" validate:"required"`
	Charger  string `json:"charger" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type DashboardUserSummary struct {
	ActiveAsset int64 `db:"active_asset"`
}

type DashboardUserData struct {
	Summary DashboardUserSummary
	Assets  []Asset
}

type Asset struct {
	Brand     string    `db:"brand" json:"brand"`
	Model     string    `db:"model" json:"model"`
	SerialNo  string    `db:"serial_no" json:"serialNo"`
	Type      string    `db:"type" json:"type"`
	Status    string    `db:"status" json:"status"`
	Owner     string    `db:"owner" json:"owner"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
type SentServiceRequest struct {
	StartDate string `json:"start_date" db:"service_start" validate:"required"`
	EndDate   string `json:"end_date" db:"service_end" validate:"required"`
}
type DashboardSummary struct {
	Total            int `json:"totalAssets" db:"total"`
	Available        int `json:"available" db:"available"`
	Assigned         int `json:"assigned" db:"assigned"`
	WaitingForRepair int `json:"waitingForRepair" db:"waiting_for_repair"`
	InService        int `json:"inService" db:"in_service"`
	Damaged          int `json:"damaged" db:"damaged"`
}
type DashboardData struct {
	Summary DashboardSummary `json:"summary"`
	Assets  []Asset          `json:"data"`
}
type UpdateAssetRequest struct {
	Brand         string `json:"brand" validate:"required"`
	Model         string `json:"model" validate:"required"`
	SerialNo      string `json:"serialNo" validate:"required"`
	Type          string `json:"type" validate:"required"`
	Owner         string `json:"owner" validate:"required"`
	WarrantyStart string `json:"warrantyStart" validate:"required"`
	WarrantyEnd   string `json:"warrantyEnd" validate:"required"`

	Laptop   *LaptopInput   `json:"laptop,omitempty"`
	Mouse    *MouseInput    `json:"mouse,omitempty"`
	Keyboard *KeyboardInput `json:"keyboard,omitempty"`
	Mobile   *MobileInput   `json:"mobile,omitempty"`
}
type AssetInfoRequest struct {
	ID     string `db:"id" json:"id"`
	Brand  string `db:"brand" json:"brand"`
	Model  string `db:"model" json:"model"`
	Status string `db:"status" json:"status"`
	Type   string `db:"type" json:"type"`
}
