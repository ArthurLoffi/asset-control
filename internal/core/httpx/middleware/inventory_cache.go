package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type capturedResponse struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (writer *capturedResponse) Write(body []byte) (int, error) {
	_, _ = writer.body.Write(body)
	return writer.ResponseWriter.Write(body)
}

func (writer *capturedResponse) WriteString(body string) (int, error) {
	_, _ = writer.body.WriteString(body)
	return writer.ResponseWriter.WriteString(body)
}

func inventoryResource(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
	if len(parts) == 0 {
		return ""
	}
	switch parts[0] {
	case "employees", "machines", "monitors", "phones", "printers", "branches", "departments":
		return parts[0]
	default:
		return ""
	}
}

// InventoryCache caches only protected, successful inventory JSON reads. The
// version changes after every successful create, update, or deactivation.
func InventoryCache(cache ResponseCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource := inventoryResource(c.Request.URL.Path)
		if resource == "" || cache == nil {
			c.Next()
			return
		}
		if c.Request.Method == http.MethodGet {
			request := c.Request.URL.Path + "?" + c.Request.URL.Query().Encode()
			body, version, hit, err := cache.Load(c.Request.Context(), resource, request)
			if err != nil {
				slog.WarnContext(c.Request.Context(), "inventory cache read", "error", err)
				c.Next()
				return
			}
			if hit {
				c.Data(http.StatusOK, "application/json; charset=utf-8", body)
				c.Abort()
				return
			}
			writer := &capturedResponse{ResponseWriter: c.Writer}
			c.Writer = writer
			c.Next()
			if c.Writer.Status() == http.StatusOK && writer.body.Len() > 0 {
				if err := cache.Store(c.Request.Context(), resource, request, version, writer.body.Bytes()); err != nil {
					slog.WarnContext(c.Request.Context(), "inventory cache write", "error", err)
				}
			}
			return
		}
		c.Next()
		if c.Request.Method != http.MethodGet &&
			c.Writer.Status() >= http.StatusOK && c.Writer.Status() < http.StatusMultipleChoices {
			invalidateCtx, cancel := context.WithTimeout(context.WithoutCancel(c.Request.Context()), 2*time.Second)
			defer cancel()
			resources := []string{resource}
			// Deactivating an employee also unassigns their equipment.
			if resource == "employees" && c.Request.Method == http.MethodDelete {
				resources = append(resources, "machines", "monitors", "phones")
			}
			for _, name := range resources {
				if err := cache.Invalidate(invalidateCtx, name); err != nil {
					slog.WarnContext(c.Request.Context(), "inventory cache invalidation", "error", err)
				}
			}
		}
	}
}
