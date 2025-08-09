package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/choipopik/gRPC-SSO/tests/suite"
	ssov1 "github.com/choipopik/protos/gen/go/sso"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "watermelon"

	passDefaultLen = 12
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	regResp, err := st.AuthClient.Register(ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: pass,
		})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	loginResp, err := st.AuthClient.Login(ctx,
		&ssov1.LoginRequest{
			Email:    email,
			Password: pass,
			AppId:    appID,
		})

	require.NoError(t, err)

	loginTime := time.Now()

	token := loginResp.GetToken()
	require.NoError(t, err)

	tokenPars, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenPars.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	// assert.Equal(t, regResp.GetUserId(), int64(claims["user_id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSecond = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSecond)
}
