package auth

import (
	"washboard/state"
	"washboard/types"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"golang.org/x/crypto/bcrypt"
)

var appState *state.Data = state.Instance()


func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 11)
	if err != nil {
		glg.Errorf("Error while hashing password: %s", err)
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		glg.Errorf("Error while comparing passwords: %s", err)
		return false
	}
	return true
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*types.User); ok {
		return jwt.MapClaims{
			types.IdentityKey: v.UserName,
		}
	}
	return jwt.MapClaims{}
}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals types.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	// eheheheheheh segur
	hashed, err := HashAndSalt([]byte(appState.Config.Password))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	if (userID == appState.Config.User && ComparePasswords(hashed, []byte(password))) {
		return &types.User{
			UserName: userID,
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &types.User{
		UserName: claims[types.IdentityKey].(string),
	}
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*types.User); ok && v.UserName == appState.Config.User {
		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
