package jwt

import (
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
	AccessTokenExpire  int64  `mapstructure:"access_token_expire"`
	RefreshTokenExpire int64  `mapstructure:"refresh_token_expire"`
	AutoLogout         int64  `mapstructure:"auto_logout"`
}

type JWTentity struct {
	ID  uuid.UUID `json:"id"` // User ID
	UID uuid.UUID `json:"uid"`
	jwt.MapClaims
}

func CreateToken(userID uuid.UUID, expire int64, secret string) (token string, uid uuid.UUID, exp int64, err error) {
	exp = time.Now().Add(time.Second * time.Duration(expire)).Unix()
	uid = uuid.New()
	claims := &JWTentity{
		ID:  userID,
		UID: uid,
		MapClaims: jwt.MapClaims{
			"exp": exp,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", uuid.Nil, 0, errors.Wrap(err, "can't create token")
	}

	return token, uid, exp, nil
}

func GenerateTokenPair(user *entity.User, accessTokenSecret, refreshTokenSecret string, accessTokenExpire, refreshTokenExpire int64) (
	cahcedToken *entity.CachedTokens,
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID uuid.UUID
	accessToken, accessUID, exp, err = CreateToken(user.ID, accessTokenExpire, accessTokenSecret)
	if err != nil {
		return nil, "", "", 0, errors.Wrap(err, "can't create access token")
	}

	refreshToken, refreshUID, _, err = CreateToken(user.ID, refreshTokenExpire, refreshTokenSecret)
	if err != nil {
		return nil, "", "", 0, errors.Wrap(err, "can't create refresh token")
	}

	cachedToken := &entity.CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	}

	return cachedToken, accessToken, refreshToken, exp, nil
}

func ValidateToken(cachedToken *entity.CachedTokens, token *JWTentity, isRefreshToken bool) error {
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

func ParseToken(tokenString string, secret string) (*JWTentity, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTentity{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't parse token")
	}

	claims, ok := token.Claims.(*JWTentity)
	if !ok {
		return nil, errors.New("can't parse token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
