package models

import "time"

type UserExist struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
}
type UserRequest struct {
	Name        string `db:"name" json:"name" valid:"required"`
	Email       string `db:"email" json:"email" valid:"email,required"`
	Role        string `db:"role" json:"role" valid:"required"`
	Type        string `db:"type" json:"type" valid:"required"`
	PhoneNumber string `db:"phone_no" json:"phone_no" valid:"required ,min =10"`
	Password    string `db:"password_hash" json:"password" valid:"required,min =6"`
	JoiningDate string `db:"joining_date" json:"joining_date" valid:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type AssignRequest struct {
	AssignedTo string `json:"assigned_to" db:"assigned_to" validate:"required"`
}
type UserInfoRequest struct {
	ID           string             `db:"id" json:"id"`
	Name         string             `db:"name" json:"name"`
	Email        string             `db:"email" json:"email"`
	PhoneNo      string             `db:"phone_no" json:"phoneNo"`
	Role         string             `db:"role" json:"role"`
	Type         string             `db:"type" json:"type"`
	CreatedAt    time.Time          `db:"created_at" json:"createdAt"`
	AssetDetails []AssetInfoRequest `json:"assetDetails"`
}
