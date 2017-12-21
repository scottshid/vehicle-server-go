package auth

import (
    "github.com/dgrijalva/jwt-go"
    "fmt"
    "net/http"
    "github.com/scottshid/app"
    "context"
)

func ValidateMiddleware(next app.AppHandler) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        tokenHeader := req.Header.Get("X-AUTH-TOKEN")
        token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("secret"), nil
        }); if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            ctx := context.WithValue(req.Context(), "username", claims["username"])
            next.ServeHTTP(w, req.WithContext(ctx))
        } else {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
        }
    })
}
