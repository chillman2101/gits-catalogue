package cache

import "github.com/chillman2101/gits-catalogue/internal/model"

// Cache keys for Book
const (
	keyBook = "book"
)

// Book methods - Get with proper typing
// This method preserves preloaded relations and supports both single and list data
func (c *CacheHelper) GetBookTypedCache(params map[string]string) (*model.BookCacheData, error) {
	var result model.BookCacheData
	err := c.GetTypedCache(keyBook, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetBookCache stores Book in cache
func (c *CacheHelper) SetBookCache(params map[string]string, data interface{}) error {
	return c.SetCache(keyBook, params, data)
}

// DeleteBookCache removes Book from cache
func (c *CacheHelper) DeleteBookCache(params map[string]string) error {
	return c.DeleteCache(keyBook, params)
}

// ScanBookCache scans Book cache keys
func (c *CacheHelper) ScanBookCache(cursor uint64) ([]string, uint64, error) {
	return c.ScanCache(keyBook+":*", cursor)
}

// InvalidateBookCache invalidates all Book cache from Redis
func (c *CacheHelper) InvalidateBookCache() error {
	cursor := uint64(0)
	for {
		keys, nextCursor, err := c.ScanBookCache(cursor)
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
