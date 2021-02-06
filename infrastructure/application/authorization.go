package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"microtecture/domain/models"
	"microtecture/infrastructure/config"
)

var (
	jwtSigningMethods = make(map[string]*jwt.SigningMethodHMAC)
)

func init() {
	jwtSigningMethods[config.HS256] = jwt.SigningMethodHS256
	jwtSigningMethods[config.HS384] = jwt.SigningMethodHS384
	jwtSigningMethods[config.HS512] = jwt.SigningMethodHS512
}

// Claims is jwt claims
type Claims struct {
	jwt.StandardClaims
	Id        uuid.UUID `json:"uuid"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Roles     []string  `json:"roles"`
}

//CreateJWT creates json web token
func (self application) CreateJWT(userid uuid.UUID, firstName string, lastName string, lifetime bool, roles ...string) (string, error) {
	var expirationTime int64
	if lifetime {
		expirationTime = time.Now().Add(time.Duration(24*365*100) * time.Hour).Unix()
	} else {
		expirationTime = time.Now().Add(time.Duration(self.Config.JWT.MaxAge) * time.Second).Unix()
	}
	claims := Claims{
		Id:        userid,
		FirstName: firstName,
		LastName:  lastName,
		Roles:     roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwtSigningMethods[self.Config.JWT.Algorithm], claims)
	tokenString, err := token.SignedString([]byte(self.Config.JWT.Secret))
	if err != nil {
		return "", errors.New(err.Error())
	}

	return tokenString, nil
}

//CreateRefreshToken creates refresh token
func (self application) CreateRefreshToken(userid uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Duration(self.Config.JWT.RefreshToken.MaxAge) * time.Second)
	claims := Claims{
		Id:    userid,
		Roles: nil,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(
		jwtSigningMethods[self.Config.JWT.RefreshToken.Algorithm], claims,
	)
	tokenString, err := token.SignedString(self.Config.JWT.RefreshToken.Secret)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return tokenString, nil
}

// Authorize checks user authorization
func (self application) Authorize(f action, roles ...string) action {
	return func(ctx *Context) error {
		tokenString, err := ctx.ReadCookie(config.ACCESS_TOKEN_NAME)
		if err != nil && err != http.ErrNoCookie {
			self.Logger.Error(fmt.Sprintf("%+v\n", err))
			return NewErrUnauthorized()
		}

		if err == http.ErrNoCookie {
			tokenString = ctx.Request.Header.Get(config.AUTHORZIATION_NAME)
			if tokenString == "" {
				return NewErrUnauthorized()
			}
		}

		keyfunc := func(token *jwt.Token) (interface{}, error) {
			return self.Config.JWT.Secret, nil
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
		if err != nil {
			self.Logger.Error(fmt.Sprintf("%+v\n", err))
			return NewErrUnauthorized()
		}

		if !token.Valid {
			roles, err := self.RefreshToken(ctx)
			if err == NewErrUnauthorized() {
				return err
			}
			if err != nil {
				self.Logger.Error(fmt.Sprintf("%+v\n", err))
				return NewErrUnauthorized()
			}
			claims.Roles = roles
		}

		doNext := false
		extBreak := false
		if len(roles) > 0 {
			for _, role := range roles {
				for _, userRole := range claims.Roles {
					if role == userRole {
						doNext = true
						extBreak = true
						break
					}
				}
				if extBreak {
					break
				}
			}
			if !doNext {
				return NewErrForbidden()
			}
		}

		return f(ctx)
	}
}

// RefreshToken refreshes token
// return roles and error
func (self application) RefreshToken(ctx *Context) ([]string, error) {
	tokenString, err := ctx.ReadCookie(config.REFRESH_TOKEN_NAME)
	if err != nil && err != http.ErrNoCookie {
		return nil, errors.New(err.Error())
	}

	if err == http.ErrNoCookie {
		tokenString = ctx.Request.Header.Get(config.REFRESH_TOKEN_NAME)
		if tokenString == "" {
			return nil, errors.New("refresh token not in cookie or header")
		}
	}

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		return self.Config.JWT.RefreshToken.Secret, nil
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if !token.Valid {
		return nil, NewErrUnauthorized()
	}

	user := models.User{Id: claims.Id}
	// TODO: get user from database

	roles := make([]string, len(user.Group.Roles))
	for i, role := range user.Group.Roles {
		roles[i] = role.EnName
	}

	jwt, err := self.CreateJWT(user.Id, user.FirstName, user.LastName, false, roles...)
	if err != nil {
		return nil, err
	}

	ctx.SetCookie(config.ACCESS_TOKEN_NAME, jwt, self.Config.JWT.MaxAge)
	ctx = ctx.WithUser(&user)

	return roles, nil
}
