package users

import "net/http"

type Handlers interface {
	// Public
	Register() http.HandlerFunc
	SignInEmail() http.HandlerFunc
	SignInUsername() http.HandlerFunc

	// Authenticated
	Me() http.HandlerFunc
	UpdateMe() http.HandlerFunc
	UpdatePasswordMe() http.HandlerFunc
	Profile() http.HandlerFunc
	UploadAvatar() http.HandlerFunc

	// Admin — Users
	Create() http.HandlerFunc
	GetMulti() http.HandlerFunc
	ListUsers() http.HandlerFunc
	Statistics() http.HandlerFunc
	NotifyList() http.HandlerFunc
	UpdateActiveStatus() http.HandlerFunc
	Get() http.HandlerFunc
	Update() http.HandlerFunc
	UpdatePassword() http.HandlerFunc
	Delete() http.HandlerFunc
	UpdateRole() http.HandlerFunc
	LogoutAllAdmin() http.HandlerFunc

	// Admin — Access Menus
	CreateAccessMenu() http.HandlerFunc
	GetAccessMenu() http.HandlerFunc
	UpdateAccessMenu() http.HandlerFunc
	DeleteAccessMenu() http.HandlerFunc
	ListAccessMenus() http.HandlerFunc

	// Admin — Files
	CreateFile() http.HandlerFunc
	GetFile() http.HandlerFunc
	UpdateFile() http.HandlerFunc
	DeleteFile() http.HandlerFunc
	ListFiles() http.HandlerFunc

	// Admin — Logs
	CreateLog() http.HandlerFunc
	GetLog() http.HandlerFunc
	ListLogs() http.HandlerFunc
	DeleteLog() http.HandlerFunc

	// Admin — Log Types
	CreateLogType() http.HandlerFunc
	GetLogType() http.HandlerFunc
	UpdateLogType() http.HandlerFunc
	DeleteLogType() http.HandlerFunc
	ListLogTypes() http.HandlerFunc

	// Admin — Menus
	CreateMenu() http.HandlerFunc
	GetMenu() http.HandlerFunc
	UpdateMenu() http.HandlerFunc
	DeleteMenu() http.HandlerFunc
	ListMenus() http.HandlerFunc
	BulkUpdateMenuSort() http.HandlerFunc

	// Admin — Modules
	CreateModule() http.HandlerFunc
	GetModule() http.HandlerFunc
	UpdateModule() http.HandlerFunc
	DeleteModule() http.HandlerFunc
	ListModules() http.HandlerFunc
	BulkUpdateModuleSort() http.HandlerFunc
}