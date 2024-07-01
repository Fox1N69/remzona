package auth

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/storage"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *logrus.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid uint64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID uint64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID uint64) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalide credentials")
)

// New returns a new instanse Auth service.
func New(
	logger *logrus.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         logger,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login - login user in system and reutnr token
// If user not found or invalidCredentials return error
func (a *Auth) Login(
	ctx context.Context,
	email, password string,
	appID uint64,
) (string, error) {
	const op = "auth.Login"

	log := a.log.WithField(op, "op")

	log.Info("auttempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")

			return "", fmt.Errorf("%w: %s", op, ErrInvalidCredentials)
		}

		a.log.Error("failde to get user", err)

		return "", fmt.Errorf("%w: %s", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credenttials", err)
		return "", fmt.Errorf("%w: %s", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%w: %s", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token")

		return "", fmt.Errorf("%w: %s", op, err)
	}

	return token, nil
}

// RegisterNewUser register new user in system and return userID
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email, pass string,
) (uint64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.WithFields(logrus.Fields{
		"op": op,
	})

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hashe: ", err)
		return 0, fmt.Errorf("%w: %s", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user: ", err)
		return 0, fmt.Errorf("%w: %s", op, err)
	}

	log.Info("user registered")

	return id, nil
}

// IsAdmin - check if the user is an admin
func (a *Auth) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.WithField(op, "op")

	log.Info("check if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %s", op, err)
	}

	log.Info("check if user is admin: ", isAdmin)

	return isAdmin, nil
}
