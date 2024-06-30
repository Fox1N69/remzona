package auth

import (
	"context"
	"sso/internal/domain/models"
	"time"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	log         *logrus.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

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

func (a *Auth) Login(
	ctx context.Context,
	email, password string,
	appID int64,
) (string, error) {

	return "", nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email, pass string,
) (int64, error) {
	panic("no implementet")
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("no implementet")
}
