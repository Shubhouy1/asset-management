# Asset Management System

A Role-Based Asset Management Backend built using Golang, Chi Router,
PostgreSQL, JWT Authentication, sqlx, and bcrypt.

This system allows organizations to manage assets, assign them to
employees, track service lifecycle, and manage users securely.

------------------------------------------------------------------------

## Tech Stack

-   Golang
-   Chi Router
-   PostgreSQL
-   sqlx
-   JWT Authentication
-   bcrypt (Password Hashing)
-   ENUM-based Database Schema

------------------------------------------------------------------------

## Features

### Authentication & Authorization

-   User Registration
-   User Login
-   JWT Token Generation
-   Session Management
-   Logout (Session Archival)
-   Role-Based Access Control (RBAC)

Supported Roles:

-   admin
-   asset-manager
-   employee
-   project-manager
-   employee-manager

------------------------------------------------------------------------

### User Management

-   Register new user
-   Login existing user
-   Logout user
-   Get all users (Admin / Asset Manager only)
-   Filter users by:
    -   name
    -   role
    -   type
    -   asset status

------------------------------------------------------------------------

### Asset Management

Supported Asset Types:

-   laptop
-   keyboard
-   mouse
-   mobile

Features:

-   Create asset with subtype details
-   Assign asset to employee
-   Send asset to service
-   Update asset (Full PUT update)
-   View assets with filters
-   Pagination support
-   View logged-in user's assigned assets

------------------------------------------------------------------------

### Asset Status Support

Supported asset statuses:

-   available
-   assigned
-   in_service
-   for_repair
-   damaged

------------------------------------------------------------------------

## API Endpoints

### Public APIs

#### Register User

POST /register

#### Login User

POST /login

------------------------------------------------------------------------

### Protected APIs (Require JWT)

All protected APIs require:

Authorization: Bearer `<JWT_TOKEN>`{=html}

------------------------------------------------------------------------

### User APIs

#### Logout

POST /logout

#### Get All Users

Access: admin, asset-manager

GET /employee

Query Parameters:

-   name
-   role
-   type
-   status

------------------------------------------------------------------------

### Asset APIs

#### Create Asset

Access: admin, asset-manager

POST /assets/

------------------------------------------------------------------------

#### Assign Asset

PUT /assets/assign/{id}

Request Body:

{ "assigned_to": "user_uuid" }

------------------------------------------------------------------------

#### Send Asset To Service

PUT /assets/sent-to-service/{id}

Request Body:

{ "startDate": "YYYY-MM-DD", "endDate": "YYYY-MM-DD" }

------------------------------------------------------------------------

#### Update Asset (Full Update)

PUT /assets/{id}

Requires full asset data including subtype object.

------------------------------------------------------------------------

#### Show Assets (Filter + Pagination)

GET /assets

Query Parameters:

-   type
-   status
-   owner
-   brand
-   model
-   serialNumber
-   page
-   limit

Default Values:

-   page = 1
-   limit = 5

------------------------------------------------------------------------

#### Get Logged-in User's Assigned Assets

GET /get-assets

Returns:

-   Active asset count
-   Asset details

------------------------------------------------------------------------

## Database Structure

### ENUM Types

-   asset_type
-   asset_status
-   user_role
-   user_type
-   owner_type
-   connection_type

------------------------------------------------------------------------

### Tables

-   users
-   assets
-   user_session
-   laptop
-   keyboard
-   mouse
-   mobile

------------------------------------------------------------------------

## Project Structure

asset-management/ 
│ 
├── router/
├── handlers/
├── middleware/
├── models/ 
├── database/ │
              └── dbhelpers/ 
├── utils/ 
└── main.go

------------------------------------------------------------------------

## Authentication Flow

1.  Register or Login
2.  Receive JWT token
3.  Use token in Authorization header
4.  Access protected routes
5.  Logout archives the session