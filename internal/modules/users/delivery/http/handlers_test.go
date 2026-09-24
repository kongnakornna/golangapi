package http

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/users"
	usersUseCase "icmongolang/internal/modules/users/usecase"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---- fakes ----

type testLogger struct{}

func (testLogger) InitLogger()                        {}
func (testLogger) Debug(args ...interface{})          {}
func (testLogger) Debugf(t string, a ...interface{})  {}
func (testLogger) Info(args ...interface{})           {}
func (testLogger) Infof(t string, a ...interface{})   {}
func (testLogger) Warn(args ...interface{})           {}
func (testLogger) Warnf(t string, a ...interface{})   {}
func (testLogger) Error(args ...interface{})          {}
func (testLogger) Errorf(t string, a ...interface{})  {}
func (testLogger) DPanic(args ...interface{})         {}
func (testLogger) DPanicf(t string, a ...interface{}) {}
func (testLogger) Fatal(args ...interface{})          {}
func (testLogger) Fatalf(t string, a ...interface{})  {}
func (testLogger) Sync() error                        { return nil }

var _ logger.Logger = testLogger{}

type mockUseCase struct {
	users.UserUseCaseI // panic if an un-stubbed method is called

	isSuper bool
	user    *models.SdUser

	list      []*models.SdUser
	listTotal int64
	stats     []users.UserStatusCount
	profile   *users.UserProfileDetail

	avatarSize int64
	avatarErr  error
}

func (m *mockUseCase) Profile(ctx context.Context, id uuid.UUID) (*users.UserProfileDetail, error) {
	return m.profile, nil
}

func (m *mockUseCase) Get(ctx context.Context, id uuid.UUID) (*models.SdUser, error) {
	return m.user, nil
}

func (m *mockUseCase) IsActive(ctx context.Context, exp models.SdUser) bool { return true }

func (m *mockUseCase) IsSuper(ctx context.Context, exp models.SdUser) bool {
	if m.user == nil {
		return m.isSuper
	}
	return m.user.IsSuperUser
}

func (m *mockUseCase) ListUsers(ctx context.Context, f users.UserFilter) ([]*models.SdUser, int64, error) {
	return m.list, m.listTotal, nil
}

func (m *mockUseCase) Statistics(ctx context.Context) ([]users.UserStatusCount, error) {
	return m.stats, nil
}

func (m *mockUseCase) UpdateAvatar(ctx context.Context, id uuid.UUID, fh *multipart.FileHeader) error {
	m.avatarSize = fh.Size
	return m.avatarErr
}

// ---- helpers ----

func rsaKeysBase64(t *testing.T) (privB64, pubB64 string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("MarshalPKIXPublicKey: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	return base64.StdEncoding.EncodeToString(privPEM), base64.StdEncoding.EncodeToString(pubPEM)
}

type testServer struct {
	router *chi.Mux
	cfg    *config.Config
	token  string
	uc     *mockUseCase
}

func newTestServer(t *testing.T, uc *mockUseCase) *testServer {
	t.Helper()
	privB64, pubB64 := rsaKeysBase64(t)
	cfg := &config.Config{
		Jwt: config.JwtConfig{
			Issuer:                     "go-dev",
			AccessTokenPublicKey:       pubB64,
			RefreshTokenPublicKey:      pubB64,
			AccessTokenExpireDuration:  3600,
			RefreshTokenExpireDuration: 1440,
		},
		Upload: config.UploadConfig{Dir: t.TempDir(), MaxAvatarBytes: 8},
	}

	mw := middleware.CreateMiddlewareManager(cfg, testLogger{}, uc)
	handler := CreateUserHandler(uc, cfg, testLogger{})
	router := chi.NewRouter()
	MapUserRoute(router, handler, mw)

	superID := uuid.New()
	if uc.user == nil {
		uc.user = &models.SdUser{ID: superID, Status: 1, IsSuperUser: true}
	}
	token, err := jwt.CreateAccessTokenRS256(uc.user.ID.String(), "u@x.com", privB64, int64(time.Hour), "go-dev")
	if err != nil {
		t.Fatalf("mint token: %v", err)
	}
	return &testServer{router: router, cfg: cfg, token: token, uc: uc}
}

func (s *testServer) do(t *testing.T, method, path string, body *bytes.Buffer, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	return rec
}

// ---- tests ----

func TestListUsersForbiddenForNonSuperUser(t *testing.T) {
	uc := &mockUseCase{
		user: &models.SdUser{ID: uuid.New(), Status: 1, IsSuperUser: false},
		list: []*models.SdUser{{}},
	}
	s := newTestServer(t, uc)

	rec := s.do(t, http.MethodGet, "/user/list", nil, "")
	if rec.Code != http.StatusForbidden {
		t.Fatalf("GET /user/list as non-superuser: code=%d body=%s, want 403", rec.Code, rec.Body.String())
	}
}

func TestStatisticsOkForSuperUser(t *testing.T) {
	uc := &mockUseCase{
		user:  &models.SdUser{ID: uuid.New(), Status: 1, IsSuperUser: true},
		stats: []users.UserStatusCount{{Status: "Active", Count: 3, RoleName: "Super Admin"}},
	}
	s := newTestServer(t, uc)

	rec := s.do(t, http.MethodGet, "/user/statistics", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /user/statistics as superuser: code=%d body=%s, want 200", rec.Code, rec.Body.String())
	}
	payload := rec.Body.String()
	if !bytes.Contains([]byte(payload), []byte(`"payload"`)) || !bytes.Contains([]byte(payload), []byte("Active")) {
		t.Fatalf("expected payload array in response, got %s", payload)
	}
}

func TestMeCreatedAtUpdatedAtFormattedInLocalTime(t *testing.T) {
	created := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 8, 31, 1, 5, 0, 0, time.UTC)
	lastSignIn := time.Date(2026, 8, 29, 23, 59, 0, 0, time.UTC)
	uc := &mockUseCase{
		user: &models.SdUser{ID: uuid.New(), Status: 1, IsSuperUser: true},
		profile: &users.UserProfileDetail{
			User:        &models.SdUser{ID: uuid.New(), CreatedDate: created, UpdatedDate: updated, Lastsignindate: lastSignIn},
			Permissions: []users.PermissionFlag{},
		},
	}
	s := newTestServer(t, uc)

	rec := s.do(t, http.MethodGet, "/user/me", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /user/me: code=%d body=%s, want 200", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		`"createddate":"2026-08-30 17:00:00"`, // UTC +7 (Asia/Bangkok)
		`"updateddate":"2026-08-31 08:05:00"`,
		`"last_sign_in":"2026-08-30 06:59:00"`,
	} {
		if !bytes.Contains([]byte(body), []byte(want)) {
			t.Fatalf("GET /user/me body missing %s, got %s", want, body)
		}
	}
}

type stubPgRepo struct {
	users.UserPgRepository
}

type stubRedisRepo struct {
	users.UserRedisRepository
}

func TestUploadAvatarTooLargeReturnsClientError(t *testing.T) {
	cfg := &config.Config{Upload: config.UploadConfig{Dir: t.TempDir(), MaxAvatarBytes: 4}}
	realUC := usersUseCase.CreateUserUseCaseI(&stubPgRepo{}, &stubRedisRepo{}, nil, cfg, testLogger{})
	handler := CreateUserHandler(realUC, cfg, testLogger{}).(*userHandler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "big.png")
	_, _ = part.Write(make([]byte, 32))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/user/me/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey,
		&models.SdUser{ID: uuid.New(), Status: 1}))

	rec := httptest.NewRecorder()
	handler.UploadAvatar()(rec, req)

	if rec.Code < 400 || rec.Code >= 500 {
		t.Fatalf("oversized avatar upload: code=%d body=%s, want 4xx client error", rec.Code, rec.Body.String())
	}
}

func TestCreatePublicRouteRemoved(t *testing.T) {
	uc := &mockUseCase{}
	s := newTestServer(t, uc)
	s.token = "" // unauthenticated request must not matter – route is gone

	rec := s.do(t, http.MethodPost, "/users", bytes.NewBufferString(`{"email":"a@b.com","password":"12345678"}`), "application/json")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("POST /users should be removed: code=%d, want 404", rec.Code)
	}
}
