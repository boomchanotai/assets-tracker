package jwt

import (
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	AccessToken  string `mapstructure:"access_token"`
	RefreshToken string `mapstructure:"refresh_token"`
	AutoLogout   int    `mapstructure:"auto_logout"`
}

func CreateToken(userID uuid.UUID, expireMinutes int, secret string) (token string, uid uuid.UUID, exp int64, err error) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()
	uid = uuid.New()
	claims := &entity.JWTentity{
		ID:  userID,
		UID: uid,
		MapClaims: jwt.MapClaims{
			"ExpiresAt": exp,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", uuid.Nil, 0, errors.Wrap(err, "can't create token")
	}

	return token, uid, exp, nil
}

func GenerateTokenPair(user *entity.User, accessTokenSecret string, refreshTokenSecret string) (
	cahcedToken *entity.CachedTokens,
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID uuid.UUID
	ExpireAccessMinutes := 15
	ExpireRefreshMinutes := 60 * 24 * 7
	accessToken, accessUID, exp, err = CreateToken(user.ID, ExpireAccessMinutes, accessTokenSecret)
	if err != nil {
		return nil, "", "", 0, errors.Wrap(err, "can't create access token")
	}

	refreshToken, refreshUID, _, err = CreateToken(user.ID, ExpireRefreshMinutes, refreshTokenSecret)
	if err != nil {
		return nil, "", "", 0, errors.Wrap(err, "can't create refresh token")
	}

	cachedToken := &entity.CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	}

	return cachedToken, accessToken, refreshToken, exp, nil
}

func ValidateToken(cachedToken *entity.CachedTokens, token *entity.JWTentity, isRefreshToken bool) error {
	var tokenUID uuid.UUID
	if isRefreshToken {
		tokenUID = cachedToken.RefreshUID
	} else {
		tokenUID = cachedToken.AccessUID
	}

	if tokenUID != token.UID {
		return errors.New("invalid token")
	}

	return nil
}

func ParseToken(tokenString string, secret string) (*entity.JWTentity, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JWTentity{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't parse token")
	}

	claims, ok := token.Claims.(*entity.JWTentity)
	if !ok {
		return nil, errors.New("can't parse token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
