package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"animex/goproxy/cors"
	"animex/goproxy/crypto"
	"animex/goproxy/pipe"
	"animex/goproxy/templates"

	"github.com/gofiber/fiber/v2"
)

var forwardedHeaders = map[string]bool{
	"content-type":     true,
	"content-length":   true,
	"content-range":    true,
	"accept-ranges":    true,
	"cache-control":    true,
	"expires":          true,
	"last-modified":    true,
	"etag":             true,
	"content-encoding": true,
	"vary":             true,
}

func ProxyHandler(c *fiber.Ctx) error {
	acao := cors.GetValidOrigin(c)

	uParam := c.Query("u")
	urlParam := c.Query("url")

	var targetUrlStr string
	if uParam != "" {
		decrypted, ok := crypto.DecryptURL(uParam)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid encrypted URL")
		}
		targetUrlStr = decrypted
	} else if urlParam != "" {
		targetUrlStr = urlParam
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Missing u parameter")
	}

	targetUrlParsed, err := url.Parse(targetUrlStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Invalid URL: %v", err))
	}

	originParam := c.Query("origin")

	req, err := http.NewRequest(http.MethodGet, targetUrlStr, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
	}

	templates.GenerateHeadersForUrl(req, targetUrlParsed, originParam)

	if headersJson := c.Query("headers"); headersJson != "" {
		var customHeaders map[string]string
		if err := json.Unmarshal([]byte(headersJson), &customHeaders); err == nil {
			for k, v := range customHeaders {
				req.Header.Set(k, v)
			}
		}
	}

	for _, h := range []string{"Range", "If-Range", "If-None-Match", "If-Modified-Since"} {
		if val := c.Get(h); val != "" {
			req.Header.Set(h, val)
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).SendString(fmt.Sprintf("Failed to fetch upstream URL: %v", err))
	}

	cors.SetCorsHeaders(c, acao)

	contentType := resp.Header.Get("Content-Type")
	contentTypeLower := strings.ToLower(contentType)

	ctIsM3u8 := strings.Contains(contentTypeLower, "mpegurl") ||
		strings.Contains(contentTypeLower, "application/vnd.apple.mpegurl") ||
		strings.Contains(contentTypeLower, "application/x-mpegurl")

	urlLooksM3u8 := strings.HasSuffix(strings.ToLower(targetUrlStr), ".m3u8")

	if ctIsM3u8 || urlLooksM3u8 {
		defer resp.Body.Close()
		m3u8Bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read upstream body")
		}

		c.Status(http.StatusOK)

		if !bytes.HasPrefix(bytes.TrimSpace(m3u8Bytes), []byte("#EXTM3U")) {
			c.Set("Content-Type", contentType)
			return c.Send(m3u8Bytes)
		}

		processed := pipe.ProcessPipeBody(m3u8Bytes, targetUrlParsed, originParam)

		c.Set("Content-Type", "application/vnd.apple.mpegurl")
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return c.Send(processed)
	}

	c.Status(resp.StatusCode)

	for key, values := range resp.Header {
		k := strings.ToLower(key)
		if forwardedHeaders[k] {
			if len(values) > 0 {
				c.Set(key, values[0])
			}
		}
	}

	if c.Get("Content-Type") == "" && contentType != "" {
		c.Set("Content-Type", contentType)
	}

	c.Response().SetBodyStream(resp.Body, int(resp.ContentLength))

	return nil
}
