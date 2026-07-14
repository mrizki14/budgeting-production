package auth

import (
	"errors"
	"strings"

	authshared "budgeting-app/golang/internal/shared/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailExists          = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrCurrentPassword      = errors.New("current password is invalid")
	ErrPasswordConfirmation = errors.New("password confirmation does not match")
	ErrPasswordTooShort     = errors.New("password must be at least 8 characters")
	ErrNameRequired         = errors.New("name is required")
)

type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return Service{repo: repo, jwtSecret: jwtSecret}
}

func (s Service) Register(name string, email string, password string) (*User, string, error) {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, "", ErrEmailExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &User{Name: name, Email: email, Password: string(hash)}
	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := authshared.CreateToken(s.jwtSecret, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s Service) Login(email string, password string) (*User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := authshared.CreateToken(s.jwtSecret, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s Service) Me(userID uint) (*User, error) {
	return s.repo.FindByID(userID)
}

func (s Service) UpdateProfile(userID uint, name string, email string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" {
		return nil, ErrNameRequired
	}

	_, err := s.repo.FindByEmailExceptID(email, userID)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Email = email
	if err := s.repo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s Service) UpdatePassword(userID uint, current string, password string, confirmation string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(current)) != nil {
		return ErrCurrentPassword
	}
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if password != confirmation {
		return ErrPasswordConfirmation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return s.repo.Save(user)
}
