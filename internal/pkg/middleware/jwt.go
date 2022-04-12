package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	authorization = "Authorization"
	RequestId = "RequestId"
)

const (
	// SigningMethodHS256 256 签名
	SigningMethodHS256 = "SIGNING_METHOD_HS256"
	// SigningMethodHS384 384 签名
	SigningMethodHS384 = "SIGNING_METHOD_HS384"
	// SigningMethodHS512 512 签名
	SigningMethodHS512 = "SIGNING_METHOD_HS512"
	// DefaultTimeout 默认超时时间2小时
	defaultTimeout = time.Hour * 2
	// DefaultIssuer 默认发行
	defaultIssuer = "CSZK"
	AuthKey = "AUTH"
)

var SecretKey = []byte("zgcszkw.com")
var (
	ErrTokenExpired = errors.New("认证失效，请重新认证")
)

func JWTAuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get(authorization)
		//requestId := ctx.Request.Header.Get(RequestId)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": ErrTokenExpired.Error()})
			ctx.Abort()
			return
		}

		authInfo, err := checkAuth(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": ErrTokenExpired.Error()})
			ctx.Abort()
			return
		}

		setAuthKey(ctx, authInfo)
	}
}

type AuthInfo struct{
	UId string `json:"uid"`
	Name string `json:"name"`
	Menus []string `json:"privileges"`
}


//JWTMiddleware 注册
type jwtMiddleware struct {
	Timeout time.Duration
	Issuer  string
	AuthInfo *AuthInfo
}

func newJWTMiddleware(authInfo *AuthInfo, timeout time.Duration) *jwtMiddleware{
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &jwtMiddleware{
		Timeout: timeout,
		Issuer: defaultIssuer,
		AuthInfo: authInfo,
	}
}

func GenerateToken(info *AuthInfo) (string, error){
	return generateToken(info, defaultTimeout)
}
func generateToken(info *AuthInfo,timeout time.Duration) (string, error){
	jwtM := newJWTMiddleware(info, timeout)
	claims := jwtM.newClaims()
	return claims.genToken()
}

func (jwtM *jwtMiddleware)newClaims() *Claims {
	return &Claims{
		jwtM.AuthInfo,
		jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(jwtM.Timeout).Unix(),
			Issuer:   jwtM.Issuer,
		},
	}
}

type Claims struct{
	Data *AuthInfo
	jwt.StandardClaims
}

func (c *Claims) genToken()(string, error) {
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(SecretKey)
}

func ParseToken(jwtToken string) (*AuthInfo,error) {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token)(interface{}, error){
		return SecretKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, err
			}
		}
		return nil, errors.New("token invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok ||  !token.Valid {
		return nil, errors.New("token invalid")
	}

	authInfo := claims.Data
	return authInfo,nil
}

func Refresh(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		claims.StandardClaims.ExpiresAt = time.Now().Add(defaultTimeout).Unix()
		return claims.genToken()
	}

	return "", errors.New("failed to refresh token")
}

func checkAuth(authorization string)(*AuthInfo, error) {
	if authorization == "" {
		return nil,ErrTokenExpired
	}
	authInfo, err := ParseToken(authorization)
	if err != nil {
		return nil, err
	}
	return authInfo, nil
}

func setAuthKey(c *gin.Context, value interface{}) {
	c.Set(AuthKey, value)
}

func GenerateRefreshToken(authInfo *AuthInfo) (string, error) {
	return generateToken(authInfo, time.Hour * 72)
}

func ParseJWT(ctx *gin.Context) (*AuthInfo, error){
	token := ctx.Request.Header.Get(authorization)
	//requestId := ctx.Request.Header.Get(RequestId)
	if token == "" {
		return nil, errors.New(ErrTokenExpired.Error())
	}

	return ParseToken(token)
}
