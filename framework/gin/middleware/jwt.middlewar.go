package middleware

//func EchoJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		authHeader := c.Request().Header.Get("Authorization")
//		tokenString := utils.ExtractToken(authHeader)
//
//		if tokenString == "" {
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid token"})
//		}
//
//		token, err := utils.ValidateToken(tokenString)
//		if err != nil || !token.Valid {
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
//		}
//
//		c.Set("user", token.Claims)
//		return next(c)
//	}
//}
