package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	CACHE_DURATION  = 2 * time.Minute
	ErrCacheDisabled = errors.New("cache not available")
)

type CacheHelper struct {
	client *redis.Client
}

func NewCacheHelper(redisClient *redis.Client) *CacheHelper {
	return &CacheHelper{client: redisClient}
}

func (c *CacheHelper) IsAvailable() bool {
	return c != nil && c.client != nil
}

func (c *CacheHelper) generateParamsHash(params map[string]string) string {
	if len(params) == 0 {
		return "empty"
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	hasher := md5.New()
	hasher.Write([]byte(strings.Join(pairs, "&")))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *CacheHelper) generateKey(baseKey string, params map[string]string) string {
	return fmt.Sprintf("%s:%s", baseKey, c.generateParamsHash(params))
}

func (c *CacheHelper) GetTypedCache(key string, params map[string]string, target interface{}) error {
	if !c.IsAvailable() {
		return ErrCacheDisabled
	}
	cacheKey := c.generateKey(key, params)
	val, err := c.client.Get(context.Background(), cacheKey).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), target)
}

func (c *CacheHelper) SetCache(key string, params map[string]string, data interface{}, ttl ...int) error {
	if !c.IsAvailable() {
		return nil
	}
	cacheKey := c.generateKey(key, params)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	duration := CACHE_DURATION
	if len(ttl) > 0 {
		duration = time.Duration(ttl[0]) * time.Hour
	}
	return c.client.Set(context.Background(), cacheKey, dataBytes, duration).Err()
}

func (c *CacheHelper) DeleteCache(key string, params map[string]string) error {
	if !c.IsAvailable() {
		return nil
	}
	cacheKey := c.generateKey(key, params)
	return c.client.Del(context.Background(), cacheKey).Err()
}

func (c *CacheHelper) ScanCache(key string, cursor uint64) ([]string, uint64, error) {
	if !c.IsAvailable() {
		return nil, 0, nil
	}
	return c.client.Scan(context.Background(), cursor, key, 0).Result()
}

func (c *CacheHelper) DeleteCacheWithoutGenerateKey(key string) error {
	if !c.IsAvailable() {
		return nil
	}
	return c.client.Del(context.Background(), key).Err()
}

func (c *CacheHelper) InvalidateCacheByPattern(pattern string) error {
	if !c.IsAvailable() {
		return nil
	}
	cursor := uint64(0)
	ctx := context.Background()
	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 0).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			if err := c.client.Del(ctx, key).Err(); err != nil {
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
