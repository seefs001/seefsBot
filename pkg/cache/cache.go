package cache

import (
	"github.com/allegro/bigcache"
	"time"
)

var Cache *bigcache.BigCache

func init() {
	Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
}
