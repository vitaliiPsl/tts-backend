package middleware

import (
	"strings"
	"vitaliiPsl/synthesizer/internal/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

func OpenRoute(jwtService *jwt.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" || !strings.Contains(authorization, "Bearer ") {
			return c.Next()
		}

		token := authorization[len("Bearer "):]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.Next()
		}

		c.Locals("userId", claims.Id)
		return c.Next()
	}
}

func ProtectedRoute(jwtService *jwt.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" || !strings.Contains(authorization, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authentication token"})
		}

		token := authorization[len("Bearer "):]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}

		c.Locals("userId", claims.Id)
		return c.Next()
	}
}
