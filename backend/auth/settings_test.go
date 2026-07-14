package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type settingsRepository struct {
	user      *User
	duplicate bool
	saved     *User
}

func (r *settingsRepository) FindByEmail(string) (*User, error) { return nil, gorm.ErrRecordNotFound }
func (r *settingsRepository) FindByID(uint) (*User, error)      { return r.user, nil }
func (r *settingsRepository) Create(*User) error                { return nil }
func (r *settingsRepository) FindByEmailExceptID(string, uint) (*User, error) {
	if r.duplicate {
		return &User{ID: 8}, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *settingsRepository) Save(user *User) error {
	r.saved = user
	return nil
}

func TestUpdateProfileRejectsEmailUsedByAnotherUser(t *testing.T) {
	repo := &settingsRepository{user: &User{ID: 3}, duplicate: true}
	service := NewService(repo, "secret")

	_, err := service.UpdateProfile(3, "Ibrahim", "used@example.com")

	if !errors.Is(err, ErrEmailExists) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}

func TestUpdateProfileSavesNameAndEmail(t *testing.T) {
	repo := &settingsRepository{user: &User{ID: 3, Name: "Old", Email: "old@example.com"}}
	service := NewService(repo, "secret")

	user, err := service.UpdateProfile(3, "Ibrahim", "new@example.com")

	if err != nil {
		t.Fatal(err)
	}
	if repo.saved == nil || user.Name != "Ibrahim" || user.Email != "new@example.com" {
		t.Fatalf("profile not saved: %#v", user)
	}
}

func TestUpdatePasswordRejectsWrongCurrentPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &settingsRepository{user: &User{ID: 3, Password: string(hash)}}
	service := NewService(repo, "secret")

	err := service.UpdatePassword(3, "wrong", "newpassword", "newpassword")

	if !errors.Is(err, ErrCurrentPassword) {
		t.Fatalf("expected current password error, got %v", err)
	}
}

func TestUpdatePasswordRejectsMismatchedConfirmation(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &settingsRepository{user: &User{ID: 3, Password: string(hash)}}
	service := NewService(repo, "secret")

	err := service.UpdatePassword(3, "password123", "newpassword", "different")

	if !errors.Is(err, ErrPasswordConfirmation) {
		t.Fatalf("expected confirmation error, got %v", err)
	}
}

func TestUpdatePasswordStoresBcryptHash(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &settingsRepository{user: &User{ID: 3, Password: string(hash)}}
	service := NewService(repo, "secret")

	err := service.UpdatePassword(3, "password123", "newpassword", "newpassword")

	if err != nil {
		t.Fatal(err)
	}
	if repo.saved == nil || bcrypt.CompareHashAndPassword([]byte(repo.saved.Password), []byte("newpassword")) != nil {
		t.Fatalf("new password was not hashed and saved")
	}
}

func TestUpdateProfileHandlerReturnsUpdatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &settingsRepository{user: &User{ID: 3, Name: "Old", Email: "old@example.com"}}
	handler := NewHandler(NewService(repo, "secret"))
	router := gin.New()
	router.PUT("/api/settings/profile", func(c *gin.Context) {
		c.Set("userID", uint(3))
		handler.UpdateProfile(c)
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/settings/profile", strings.NewReader(`{"name":"Ibrahim","email":"new@example.com"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), "Profil berhasil diperbarui.") {
		t.Fatalf("unexpected response %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestUpdatePasswordHandlerReturnsFieldError(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(&settingsRepository{user: &User{ID: 3, Password: string(hash)}}, "secret"))
	router := gin.New()
	router.PUT("/api/settings/password", func(c *gin.Context) {
		c.Set("userID", uint(3))
		handler.UpdatePassword(c)
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/settings/password", strings.NewReader(`{"current_password":"wrong","password":"newpassword","password_confirmation":"newpassword"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity || !strings.Contains(recorder.Body.String(), "current_password") {
		t.Fatalf("unexpected response %d: %s", recorder.Code, recorder.Body.String())
	}
}
