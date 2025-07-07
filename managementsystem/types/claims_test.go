package types

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	secret     = []byte("secret")
	username   = "test_username"
	password   = "test_password"
	id         = "test_id"
	permission = "read"
)

func TestJwtCustomClaims_IsAccess(t *testing.T) {
	tests := []struct {
		name  string
		claim JwtCustomClaims
		want  bool
	}{
		{
			name:  "Is Access",
			claim: JwtCustomClaims{Types: REFRESHTYPE},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.claim.IsAccess(); got != tt.want {
				t.Errorf("JwtCustomClaims.IsAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJwtCustomClaims_IsRefresh(t *testing.T) {
	tests := []struct {
		name  string
		claim JwtCustomClaims
		want  bool
	}{
		{
			name:  "Is Refresh",
			claim: JwtCustomClaims{Types: REFRESHTYPE},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.claim.IsRefresh(); got != tt.want {
				t.Errorf("JwtCustomClaims.IsRefresh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	t.Run("Generate Token Success", func(t *testing.T) {
		claims, err := GenerateToken(secret, username, password, id)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, password, claims.Password)
		assert.Equal(t, id, claims.ClaimID)
	})

	t.Run("Generate Token with nil secret", func(t *testing.T) {
		claims, err := GenerateToken(nil, username, password, id)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestJwtCustomClaims_AccessToken(t *testing.T) {
	t.Run("Generate Access Token Success", func(t *testing.T) {
		claims, _ := GenerateToken(secret, username, password, id)
		token, err := claims.AccessToken()
		assert.NoError(t, err)
		assert.NotNil(t, token)

		parsedClaims := &JwtCustomClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		assert.NoError(t, err)
		assert.NotNil(t, parsedToken)
		assert.Equal(t, claims.Username, parsedClaims.Username)
		assert.Equal(t, claims.Password, parsedClaims.Password)
		assert.Equal(t, claims.ClaimID, parsedClaims.ClaimID)
		assert.Equal(t, claims.Types, parsedClaims.Types)
	})

	// t.Run("Access Token with nil secret", func(t *testing.T) {
	// 	claims, _ := GenerateToken([]byte{}, username, password, id)
	// 	token, err := claims.AccessToken()
	// 	assert.Error(t, err)
	// 	assert.Empty(t, token)
	// })
}

func TestJwtCustomClaims_RefreshToken(t *testing.T) {
	t.Run("Generate Refresh Token Success", func(t *testing.T) {
		claims, _ := GenerateToken(secret, username, password, id)
		token, err := claims.RefreshToken()
		assert.NoError(t, err)
		assert.NotNil(t, token)

		parseClaims := &JwtCustomClaims{}
		parseToken, err := jwt.ParseWithClaims(token, parseClaims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		assert.NoError(t, err)
		assert.NotNil(t, parseToken)
		assert.Equal(t, claims.Username, parseClaims.Username)
		assert.Equal(t, claims.Password, parseClaims.Password)
		assert.Equal(t, claims.ClaimID, parseClaims.ClaimID)
		assert.Equal(t, claims.Types, parseClaims.Types)
	})

	// t.Run("Refresh Token with nil secret", func(t *testing.T) {
	// 	claims, _ := GenerateToken(secret, username, password, id)
	// 	token, err := claims.RefreshToken()
	// 	assert.Error(t, err)
	// 	assert.Empty(t, token)
	// })
}

func TestRestricted(t *testing.T) {
	cliams, _ := GenerateToken(secret, username, password, id)
	accessToken, _ := cliams.AccessToken()

	token, _ := jwt.ParseWithClaims(accessToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	t.Run("Valid Token", func(t *testing.T) {
		e := echo.New()
		req := e.NewContext(nil, nil)
		req.Set(Authorization, token)

		cliams, ok := Restricted(req)
		assert.True(t, ok)
		assert.NotNil(t, cliams)
		assert.Equal(t, username, cliams.Username)
		assert.Equal(t, password, cliams.Password)
		assert.Equal(t, id, cliams.ClaimID)
		assert.Equal(t, ACCESSTYPE, cliams.Types)
	})
}

func TestParseWithClaims(t *testing.T) {
	t.Run("Parse With Claims Success", func(t *testing.T) {
		cliams, _ := GenerateToken(secret, username, password, id)
		accesToken, _ := cliams.AccessToken()
		ParseClaims, err := ParseWithClaims(secret, accesToken)
		assert.NoError(t, err)
		assert.NotNil(t, ParseClaims)
		assert.Equal(t, username, ParseClaims.Username)
		assert.Equal(t, password, ParseClaims.Password)
		assert.Equal(t, id, ParseClaims.ClaimID)
		assert.Equal(t, ACCESSTYPE, ParseClaims.Types)
	})

	t.Run("Parse with valid claims", func(t *testing.T) {
		_, err := ParseWithClaims(secret, "")
		assert.Error(t, err)
	})

	t.Run("Parse with wrong secret", func(t *testing.T) {
		claims, _ := GenerateToken([]byte{}, username, password, id)
		accessToken, _ := claims.AccessToken()
		_, err := ParseWithClaims([]byte("wrongsecret"), accessToken)
		assert.Error(t, err)
	})
}
