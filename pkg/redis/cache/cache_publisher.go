package cache

import "github.com/chillman2101/gits-catalogue/internal/model"

// Cache keys for Publisher
const (
	keyPublisher = "publisher"
)

// Publisher methods - Get with proper typing
// This method preserves preloaded relations and supports both single and list data
func (c *CacheHelper) GetPublisherTypedCache(params map[string]string) (*model.PublisherCacheData, error) {
	var result model.PublisherCacheData
	err := c.GetTypedCache(keyPublisher, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetPublisherCache stores Publisher in cache
func (c *CacheHelper) SetPublisherCache(params map[string]string, data interface{}) error {
	return c.SetCache(keyPublisher, params, data)
}

// DeletePublisherCache removes Publisher from cache
func (c *CacheHelper) DeletePublisherCache(params map[string]string) error {
	return c.DeleteCache(keyPublisher, params)
}

// ScanPublisherCache scans Publisher cache keys
func (c *CacheHelper) ScanPublisherCache(cursor uint64) ([]string, uint64, error) {
	return c.ScanCache(keyPublisher+":*", cursor)
}

// InvalidatePublisherCache invalidates all Publisher cache from Redis
func (c *CacheHelper) InvalidatePublisherCache() error {
	cursor := uint64(0)
	for {
		keys, nextCursor, err := c.ScanPublisherCache(cursor)
		if err != nil {
			return err
		}

		// Delete each key
		for _, key := range keys {
			if err := c.DeleteCacheWithoutGenerateKey(key); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
