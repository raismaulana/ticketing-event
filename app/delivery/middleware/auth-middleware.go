package middleware

import (
	"log"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

func AuthMiddleware(jwtCase usecase.JWTCase, e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := jwtCase.ValidateToken(tokenString)
		if !token.Valid {
			log.Println(err)
			response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		ok := checkPermission(claims, c.Request, e)
		if !ok {
			response := helper.BuildErrorResponse("Forbidden", "You don't have authorization to access this URL", nil)
			c.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		c.Next()
	}

}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func checkPermission(claims jwt.MapClaims, r *http.Request, e *casbin.Enforcer) bool {
	_, ok := claims["user_id"].(string)
	role, ok2 := claims["role"].(string)
	if !(ok && ok2) {
		return false
	}
	user := role
	method := r.Method
	path := r.URL.Path
	ok3 := e.Enforce(user, path, method)
	return ok3
}
