package usecase

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/users"
	"icmongolang/pkg/httpErrors"

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

type mockPgRepo struct {
	users.UserPgRepository // panic if an un-stubbed method is called

	user        *models.SdUser
	profile     *users.UserProfileDetail
	profileHits int

	listUsers []*models.SdUser
	listTotal int64
	listErr   error

	stats []users.UserStatusCount

	notified      []*models.SdUser
	notifiedChann users.NotificationChannel

	activeStatusID  uuid.UUID
	activeStatusVal int16

	avatarID   uuid.UUID
	avatarPath string
	avatarFile string

	getByEmailUser *models.SdUser
}

func (m *mockPgRepo) Get(ctx context.Context, id uuid.UUID) (*models.SdUser, error) {
	return m.user, nil
}

func (m *mockPgRepo) GetProfileWithPermissions(ctx context.Context, id uuid.UUID) (*users.UserProfileDetail, error) {
	m.profileHits++
	if m.profile != nil {
		if m.profile.User == nil {
			m.profile.User = &models.SdUser{}
		}
		m.profile.User.ID = id
	}
	return m.profile, nil
}

func (m *mockPgRepo) ListUsers(ctx context.Context, f users.UserFilter) ([]*models.SdUser, int64, error) {
	return m.listUsers, m.listTotal, m.listErr
}

func (m *mockPgRepo) GetStatistics(ctx context.Context) ([]users.UserStatusCount, error) {
	return m.stats, nil
}

func (m *mockPgRepo) GetActiveByNotification(ctx context.Context, ch users.NotificationChannel) ([]*models.SdUser, error) {
	m.notifiedChann = ch
	return m.notified, nil
}

func (m *mockPgRepo) UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error {
	m.activeStatusID = id
	m.activeStatusVal = status
	return nil
}

func (m *mockPgRepo) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarPath, avatar string) error {
	m.avatarID = id
	m.avatarPath = avatarPath
	m.avatarFile = avatar
	return nil
}

func (m *mockPgRepo) GetByEmail(ctx context.Context, identifier string) (*models.SdUser, error) {
	if m.getByEmailUser == nil {
		return nil, errors.New("not found")
	}
	return m.getByEmailUser, nil
}

type mockRedisRepo struct {
	users.UserRedisRepository // panic if an un-stubbed method is called

	detail     *users.UserProfileDetail
	setCalled  bool
	setTTL     int
	deletes    []string
	sadds      []string
	expiresKey string
	expiresSec int
}

func (m *mockRedisRepo) Get(ctx context.Context, key string) (*models.SdUser, error) {
	return nil, nil
}

func (m *mockRedisRepo) Create(ctx context.Context, key string, exp *models.SdUser, seconds int) error {
	return nil
}

func (m *mockRedisRepo) Delete(ctx context.Context, key string) error {
	m.deletes = append(m.deletes, key)
	return nil
}

func (m *mockRedisRepo) GetDetail(ctx context.Context, key string) (*users.UserProfileDetail, error) {
	return m.detail, nil
}

func (m *mockRedisRepo) SetDetail(ctx context.Context, key string, d *users.UserProfileDetail, seconds int) error {
	m.setCalled = true
	m.setTTL = seconds
	return nil
}

func (m *mockRedisRepo) Sadd(ctx context.Context, key string, value string) error {
	m.sadds = append(m.sadds, key)
	return nil
}

func (m *mockRedisRepo) Expire(ctx context.Context, key string, seconds int) error {
	m.expiresKey = key
	m.expiresSec = seconds
	return nil
}

// ---- helpers ----

func newTestUseCase(pg *mockPgRepo, rd *mockRedisRepo, cfg *config.Config) users.UserUseCaseI {
	if cfg == nil {
		cfg = &config.Config{}
	}
	return CreateUserUseCaseI(pg, rd, nil, cfg, testLogger{})
}

// makeFileHeader builds a real *multipart.FileHeader with working Open().
func makeFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		t.Fatalf("ParseMultipartForm: %v", err)
	}
	_, fh, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("FormFile: %v", err)
	}
	return fh
}

func restMsg(err error) string {
	var rest httpErrors.ErrRest
	if errors.As(err, &rest) {
		return rest.GetMsg()
	}
	return err.Error()
}

// ---- tests ----

func TestProfileCacheHitSkipsDB(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{profile: &users.UserProfileDetail{RoleTitle: "Admin"}}
	rd := &mockRedisRepo{detail: &users.UserProfileDetail{RoleTitle: "Cached"}}
	uc := newTestUseCase(pg, rd, nil)

	detail, err := uc.Profile(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if detail.RoleTitle != "Cached" {
		t.Fatalf("expected cached profile, got %q", detail.RoleTitle)
	}
	if pg.profileHits != 0 {
		t.Fatalf("DB should not be hit on cache hit, hits=%d", pg.profileHits)
	}
	if rd.setCalled {
		t.Fatal("SetDetail should not be called on cache hit")
	}
}

func TestProfileCacheMissLoadsAndCaches(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{profile: &users.UserProfileDetail{RoleTitle: "Admin"}}
	rd := &mockRedisRepo{}
	uc := newTestUseCase(pg, rd, nil)

	detail, err := uc.Profile(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if detail == nil || detail.RoleTitle != "Admin" {
		t.Fatalf("unexpected detail: %+v", detail)
	}
	if pg.profileHits != 1 {
		t.Fatalf("repo should be called once, hits=%d", pg.profileHits)
	}
	if !rd.setCalled || rd.setTTL != 3600 {
		t.Fatalf("expected SetDetail TTL=3600, called=%v ttl=%d", rd.setCalled, rd.setTTL)
	}
}

func TestUpdateActiveStatusValidation(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{user: &models.SdUser{}}
	rd := &mockRedisRepo{}
	uc := newTestUseCase(pg, rd, nil)

	if err := uc.UpdateActiveStatus(context.Background(), id, 2); err == nil {
		t.Fatal("expected validation error for status=2")
	}
	if pg.activeStatusID == id {
		t.Fatal("repo must not be called for invalid status")
	}
}

func TestUpdateActiveStatusSuccessInvalidatesCache(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{user: &models.SdUser{ID: id}}
	rd := &mockRedisRepo{}
	uc := newTestUseCase(pg, rd, nil)

	if err := uc.UpdateActiveStatus(context.Background(), id, 1); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if pg.activeStatusID != id || pg.activeStatusVal != 1 {
		t.Fatalf("repo args mismatch: id=%s status=%d", pg.activeStatusID, pg.activeStatusVal)
	}
	wantDeletes := map[string]bool{
		uc.GenerateRedisUserKey(id):    false,
		uc.GenerateRedisProfileKey(id): false,
	}
	for _, k := range rd.deletes {
		if _, ok := wantDeletes[k]; ok {
			wantDeletes[k] = true
		}
	}
	for k, seen := range wantDeletes {
		if !seen {
			t.Fatalf("cache key %q was not invalidated (deletes=%v)", k, rd.deletes)
		}
	}
}

func TestActiveByNotification(t *testing.T) {
	pg := &mockPgRepo{user: &models.SdUser{}, notified: []*models.SdUser{{}}}
	rd := &mockRedisRepo{}
	uc := newTestUseCase(pg, rd, nil)

	if _, err := uc.ActiveByNotification(context.Background(), users.NotificationChannel("web")); err == nil {
		t.Fatal("expected validation error for channel=web")
	}

	out, err := uc.ActiveByNotification(context.Background(), users.NotifyEmail)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if pg.notifiedChann != users.NotifyEmail {
		t.Fatalf("expected pass-through channel=email, got %s", pg.notifiedChann)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 user, got %d", len(out))
	}
}

func TestUpdateAvatarTooLarge(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{}
	rd := &mockRedisRepo{}
	cfg := &config.Config{Upload: config.UploadConfig{MaxAvatarBytes: 4}}
	uc := newTestUseCase(pg, rd, cfg)

	fh := makeFileHeader(t, "avatar.png", []byte("0123456789"))
	err := uc.UpdateAvatar(context.Background(), id, fh)
	if err == nil {
		t.Fatal("expected size error")
	}
	if !strings.Contains(restMsg(err), "max size") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestUpdateAvatarInvalidExtension(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{}
	rd := &mockRedisRepo{}
	uc := newTestUseCase(pg, rd, nil)

	fh := makeFileHeader(t, "avatar.gif", []byte("gif89a"))
	if err := uc.UpdateAvatar(context.Background(), id, fh); err == nil {
		t.Fatal("expected extension error for .gif")
	}
	if pg.avatarFile != "" {
		t.Fatal("pgRepo.UpdateAvatar must not be called for invalid ext")
	}
}

func TestUpdateAvatarSuccess(t *testing.T) {
	id := uuid.New()
	dir := t.TempDir()
	pg := &mockPgRepo{user: &models.SdUser{ID: id}}
	rd := &mockRedisRepo{}
	cfg := &config.Config{Upload: config.UploadConfig{Dir: dir, MaxAvatarBytes: 2 << 20}}
	uc := newTestUseCase(pg, rd, cfg)

	fh := makeFileHeader(t, "me.PNG", []byte("fakepng"))
	if err := uc.UpdateAvatar(context.Background(), id, fh); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	wantPath := "/uploads/avatar/" + id.String() + ".png"
	if pg.avatarPath != wantPath {
		t.Fatalf("avatarPath = %q, want %q", pg.avatarPath, wantPath)
	}
	saved := filepath.Join(dir, "avatar", id.String()+".png")
	if _, err := os.Stat(saved); err != nil {
		t.Fatalf("file not saved to disk: %v", err)
	}
	found := false
	for _, k := range rd.deletes {
		if k == uc.GenerateRedisProfileKey(id) {
			found = true
		}
	}
	if !found {
		t.Fatalf("profile cache not invalidated (deletes=%v)", rd.deletes)
	}
}

func TestSignInByUsernameEnumerationSafe(t *testing.T) {
	pg := &mockPgRepo{} // GetByEmail returns error -> user not found
	rd := &mockRedisRepo{}
	uc := newTestUseCase(pg, rd, nil)

	_, _, errByUsername := uc.SignInByUsername(context.Background(), "ghost", "pw123456")
	_, _, errByEmail := uc.SignIn(context.Background(), "ghost@x.com", "pw123456")

	msgByUsername := restMsg(errByUsername)
	msgByEmail := restMsg(errByEmail)
	if msgByUsername != "invalid credentials" {
		t.Fatalf("SignInByUsername msg = %q, want \"invalid credentials\"", msgByUsername)
	}
	if msgByUsername != msgByEmail {
		t.Fatalf("SignIn (%q) and SignInByUsername (%q) messages must match", msgByEmail, msgByUsername)
	}
}

func TestStoreRefreshTokenAppliesTTL(t *testing.T) {
	id := uuid.New()
	pg := &mockPgRepo{}
	rd := &mockRedisRepo{}
	cfg := &config.Config{Jwt: config.JwtConfig{RefreshTokenExpireDuration: 1440}}
	ucImpl := newTestUseCase(pg, rd, cfg).(*userUseCase)

	if err := ucImpl.storeRefreshToken(context.Background(), id, "refresh-token"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(rd.sadds) != 1 || rd.sadds[0] != ucImpl.GenerateRedisRefreshTokenKey(id) {
		t.Fatalf("Sadd keys = %v, want [%s]", rd.sadds, ucImpl.GenerateRedisRefreshTokenKey(id))
	}
	if rd.expiresKey != ucImpl.GenerateRedisRefreshTokenKey(id) || rd.expiresSec != 1440*60 {
		t.Fatalf("Expire key=%s sec=%d, want key=%s sec=%d",
			rd.expiresKey, rd.expiresSec, ucImpl.GenerateRedisRefreshTokenKey(id), 1440*60)
	}
}
