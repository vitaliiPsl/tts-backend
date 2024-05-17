package auth

import (
	"strings"
	"vitaliiPsl/synthesizer/internal/auth/jwt"
	"vitaliiPsl/synthesizer/internal/user"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtService  *jwt.JwtService
	userService *user.UserService
}

func NewAuthMiddleware(jwtService *jwt.JwtService, userService *user.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		userService: userService,
	}
}

func (m *AuthMiddleware) OpenRoute() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" || !strings.Contains(authorization, "Bearer ") {
			return c.Next()
		}

		token := authorization[len("Bearer "):]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			return c.Next()
		}

		return m.fetchUser(c, claims.Id)
	}
}

func (m *AuthMiddleware) ProtectedRoute(roles ...user.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" || !strings.Contains(authorization, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authentication token"})
		}

		token := authorization[len("Bearer "):]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}

		return m.fetchUser(c, claims.Id, roles...)
	}
}

func (m *AuthMiddleware) fetchUser(c *fiber.Ctx, id string, roles ...user.UserRole) error {
	userDto, err := m.userService.FindById(id)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	if userDto.Status != user.StatusActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User is not active"})
	}

	if !m.checkUserRole(userDto, roles...) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not enough permissions"})
	}

	c.Locals("user", userDto)
	return c.Next()
}

func (m *AuthMiddleware) checkUserRole(userDto *user.UserDto, roles ...user.UserRole) bool {
	if len(roles) == 0 {
		return true
	}

	for _, role := range roles {
		if userDto.Role == role {
			return true
		}
	}

	return false
}
