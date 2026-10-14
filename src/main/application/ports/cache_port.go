package ports

import "context"

// InventoryCache caches successful inventory reads by resource and request URL.
// A version returned by Load keeps a concurrent write from publishing an old read
// under the new version.
type InventoryCache interface {
	Load(ctx context.Context, resource, request string) (body []byte, version string, hit bool, err error)
	Store(ctx context.Context, resource, request, version string, body []byte) error
	Invalidate(ctx context.Context, resource string) error
}
