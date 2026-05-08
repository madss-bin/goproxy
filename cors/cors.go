package cors

import (
	"strings"

	"animex/goproxy/config"

	"github.com/gofiber/fiber/v2"
)

func GetValidOrigin(c *fiber.Ctx) string {
	if !config.ENABLE_CORS {
		return ""
	}

	origin := c.Get("Origin")
	if origin != "" {
		for _, o := range config.ALLOWED_ORIGIN {
			if o == origin {
				return origin
			}
		}
	}

	referer := c.Get("Referer")
	if referer != "" {
		for _, o := range config.ALLOWED_ORIGIN {
			if strings.HasPrefix(referer, o) {
				return o
			}
		}
	}

	return ""
}

func SetCorsHeaders(c *fiber.Ctx, origin string) {
	acao := origin
	if acao == "" {
		acao = "*"
	}
	c.Set("Access-Control-Allow-Origin", acao)
	c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
	c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Range, X-Requested-With, Origin, Accept, Accept-Encoding, Accept-Language, Cache-Control, Pragma, Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, Connection")
	c.Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges, Content-Type, Cache-Control, Expires, Vary, ETag, Last-Modified")
	c.Set("Cross-Origin-Resource-Policy", "cross-origin")
	c.Set("Vary", "Origin")
}

func HandleOptions(c *fiber.Ctx) error {
	origin := GetValidOrigin(c)
	if config.ENABLE_CORS && origin == "" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	SetCorsHeaders(c, origin)
	c.Set("Access-Control-Max-Age", "86400")
	return c.SendStatus(fiber.StatusOK)
}
