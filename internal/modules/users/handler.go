package users

import "net/http"

type Handlers interface {
	// Public
	Register() http.HandlerFunc
	SignInEmail() http.HandlerFunc    // POST /api/auth/signin (email)
	SignInUsername() http.HandlerFunc // POST /api/auth/login (username)

	// Authenticated
	Me() http.HandlerFunc
	UpdateMe() http.HandlerFunc
	UpdatePasswordMe() http.HandlerFunc
	Profile() http.HandlerFunc
	UploadAvatar() http.HandlerFunc

	// Admin
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
}
