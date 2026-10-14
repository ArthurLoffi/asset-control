package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeInventoryCache struct {
	versions map[string]int
	entries  map[string][]byte
}

func (cache *fakeInventoryCache) key(resource, request, version string) string {
	return resource + ":" + request + ":" + version
}

func (cache *fakeInventoryCache) Load(_ context.Context, resource, request string) ([]byte, string, bool, error) {
	version := fmt.Sprint(cache.versions[resource])
	body, hit := cache.entries[cache.key(resource, request, version)]
	return body, version, hit, nil
}

func (cache *fakeInventoryCache) Store(_ context.Context, resource, request, version string, body []byte) error {
	cache.entries[cache.key(resource, request, version)] = append([]byte(nil), body...)
	return nil
}

func (cache *fakeInventoryCache) Invalidate(_ context.Context, resource string) error {
	cache.versions[resource]++
	return nil
}

func TestInventoryCacheReadAndMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &fakeInventoryCache{versions: map[string]int{}, entries: map[string][]byte{}}
	router := gin.New()
	router.Use(InventoryCache(cache))
	reads := 0
	router.GET("/api/machines/id", func(c *gin.Context) {
		reads++
		c.JSON(http.StatusOK, gin.H{"data": reads})
	})
	router.PUT("/api/machines", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/api/machines/deactivate", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/api/machines", func(c *gin.Context) { c.Status(http.StatusBadRequest) })

	call := func(method, target string) string {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(method, target, nil))
		return recorder.Body.String()
	}
	first := call(http.MethodGet, "/api/machines/id?id=1")
	if second := call(http.MethodGet, "/api/machines/id?id=1"); second != first || reads != 1 {
		t.Fatalf("cache miss on repeated read: %q, %q, reads %d", first, second, reads)
	}
	call(http.MethodPost, "/api/machines")
	call(http.MethodGet, "/api/machines/id?id=1")
	if reads != 1 {
		t.Fatalf("failed mutation invalidated cache: %d reads", reads)
	}
	call(http.MethodPut, "/api/machines?id=1")
	call(http.MethodGet, "/api/machines/id?id=1")
	if reads != 2 {
		t.Fatalf("update did not invalidate cache: %d reads", reads)
	}
	call(http.MethodPost, "/api/machines/deactivate?id=1")
	call(http.MethodGet, "/api/machines/id?id=1")
	if reads != 3 {
		t.Fatalf("deactivation did not invalidate cache: %d reads", reads)
	}
}

func TestInventoryCacheDoesNotCacheEmails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &fakeInventoryCache{versions: map[string]int{}, entries: map[string][]byte{}}
	router := gin.New()
	router.Use(InventoryCache(cache))
	reads := 0
	router.GET("/api/emails", func(c *gin.Context) {
		reads++
		c.JSON(http.StatusOK, gin.H{"password_hash": "secret"})
	})
	for range 2 {
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/emails", nil))
	}
	if reads != 2 || len(cache.entries) != 0 {
		t.Fatalf("email response was cached: %d reads, %d entries", reads, len(cache.entries))
	}
}

func TestEmployeeDeactivationInvalidatesAssignedEquipment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &fakeInventoryCache{versions: map[string]int{}, entries: map[string][]byte{}}
	router := gin.New()
	router.Use(InventoryCache(cache))
	router.POST("/api/employees/deactivate", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/employees/deactivate?id=1", nil))
	for _, resource := range []string{"employees", "machines", "monitors", "phones"} {
		if cache.versions[resource] != 1 {
			t.Fatalf("%s cache version was not advanced", resource)
		}
	}
}
