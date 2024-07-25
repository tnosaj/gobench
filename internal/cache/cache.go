package cache

import (
	"strings"

	"github.com/samborkent/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

type CacheValues interface {
	Load() error
	Save() error
	GetRandom() (uuid.UUID, error)
	Put(uuid.UUID) error
}

func NewCache(cache string, randomizer internal.Random) CacheValues {
	cacheParts := strings.Split(cache, ":")
	switch cacheParts[0] {
	case "none":
		return newNullCache()
	case "redis":
		if len(cacheParts) == 2 {
			return newRedisCache(cacheParts[1], "6379", randomizer)
		} else if len(cacheParts) == 3 {
			return newRedisCache(cacheParts[1], cacheParts[2], randomizer)
		}
		logrus.Fatalf("Could not split redis cache into parts: %s", cache)
		return nil
	case "file":
		if len(cacheParts) != 2 {
			logrus.Fatalf("Could not split file cache into 2 parts: %s", cache)
		}
		return newFileCache(cacheParts[1], randomizer)
	case "memory":
		if len(cacheParts) != 2 {
			logrus.Fatalf("Could not split file cache into 2 parts: %s", cache)
		}
		return newMemoryCache(cacheParts[1], randomizer)
	default:
		logrus.Fatalf("unknown cachetype: %s", cacheParts[0])
		return nil
	}
}
