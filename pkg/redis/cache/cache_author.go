package cache

import "github.com/chillman2101/gits-catalogue/internal/model"

// Cache keys for Author
const (
	keyAuthor = "author"
)

// Author methods - Get with proper typing
// This method preserves preloaded relations and supports both single and list data
func (c *CacheHelper) GetAuthorTypedCache(params map[string]string) (*model.AuthorCacheData, error) {
	var result model.AuthorCacheData
	err := c.GetTypedCache(keyAuthor, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAuthorCache stores Author in cache
func (c *CacheHelper) SetAuthorCache(params map[string]string, data interface{}) error {
	return c.SetCache(keyAuthor, params, data)
}

// DeleteAuthorCache removes Author from cache
func (c *CacheHelper) DeleteAuthorCache(params map[string]string) error {
	return c.DeleteCache(keyAuthor, params)
}

// ScanAuthorCache scans Author cache keys
func (c *CacheHelper) ScanAuthorCache(cursor uint64) ([]string, uint64, error) {
	return c.ScanCache(keyAuthor+":*", cursor)
}

// InvalidateAuthorCache invalidates all Author cache from Redis
func (c *CacheHelper) InvalidateAuthorCache() error {
	cursor := uint64(0)
	for {
		keys, nextCursor, err := c.ScanAuthorCache(cursor)
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
