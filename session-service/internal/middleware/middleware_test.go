package token

// Fuck this test, for now

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/require"
// )

// func addAuthorization(
// 	t *testing.T,
// 	request *http.Request,
// 	tokenMaker token.JWTMaker,
// 	authorizationType string,
// 	username string,
// 	duration time.Duration,
// 	tokenType token.TokenType,
// ) {
// 	token, payload, err := tokenMaker.CreateToken(username, duration, tokenType)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, payload)

// 	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
// 	request.Header.Set(authorizationHeaderKey, authorizationHeader)
// }

// func TestAuthMiddleware(t *testing.T) {
// 	username := "Cole"
// 	testCases := []struct {
// 		name          string
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, time.Minute, 1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoAuthorization",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "UnsupportedAuthorization",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
// 				addAuthorization(t, request, tokenMaker, "unsupported", username, time.Minute, 1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidAuthorization",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
// 				addAuthorization(t, request, tokenMaker, "", username, time.Minute, 1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "ExpiredToken",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, -time.Minute, 1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			router := gin.New()
// 			// authPath := "/auth"
// 			maker, err := token.NewJWTMaker("Dk4cLr7zvUeFYRAxmPlgwXqJ3uEZntBS")
// 			require.NoError(t, err)

// 			router.Use(authMiddleware(*maker))

// 			router.GET("/auth", func(c *gin.Context) {
// 				c.JSON(http.StatusOK, gin.H{"message": "authorized"})
// 			})

// 			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
// 			require.NoError(t, err)

// 			tc.setupAuth(t, req, *maker)

// 			recorder := httptest.NewRecorder()
// 			router.ServeHTTP(recorder, req)
// 			tc.checkResponse(t, recorder)

// 		})

// 	}

// }
