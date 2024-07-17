package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/samborkent/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

type RedisCache struct {
	Redis      *redis.Client
	Randomizer internal.Random
	Channel    chan uuid.UUID
	CacheSize  int
}

func newRedisCache(url, port string, randomizer internal.Random) *RedisCache {
	logrus.Infof("Starting new redis cache for %s", url)

	channel := make(chan uuid.UUID)

	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", url, port),
	})
	_, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(10), // hardcoded
	)
	defer cancel()
	res := cli.Ping().Err()
	if res != nil {
		logrus.Fatalf("Error connecting to redis: %s", res)
	} else {
		logrus.Debug("Successfully connected to redis")
	}

	return &RedisCache{Redis: cli, Randomizer: randomizer, Channel: channel}
}

func (rc *RedisCache) GetRandom() (uuid.UUID, error) {
	logrus.Trace("rediscache get random")

	uid, err := rc.Redis.Get(
		strconv.Itoa(rc.Randomizer.Intn(rc.CacheSize)),
	).Result()
	if err == redis.Nil {
		logrus.Infof("key does not exist")
		return uuid.NewV7(), nil
	} else if err != nil {
		return uuid.NewV7(), fmt.Errorf("Error getting key: %s", err)
	}
	id, err := uuid.StringToUUID(uid)
	if err != nil {
		return uuid.NewV7(), fmt.Errorf("Error converting key: %s", err)
	}
	return id, nil
}

func (rc *RedisCache) Put(uuid uuid.UUID) error {
	logrus.Tracef("rediscache put %s", uuid)
	rc.Channel <- uuid
	return nil
}
func (rc *RedisCache) asyncPut(c chan uuid.UUID) {
	var err error
	for uuid := range c {

		rc.CacheSize++
		err = rc.Redis.Set(strconv.Itoa(rc.CacheSize), uuid.String(), 0).Err()
		if err != nil {
			logrus.Errorf("rediscache could not set new cache key: %s", err)
		}
		err = rc.Redis.Set("cacheSize", strconv.Itoa(rc.CacheSize), 0).Err()
		if err != nil {
			logrus.Errorf("rediscache could not increment cacheSize: %s", err)
		}
		err = nil
	}
	logrus.Info("Channel Closed, shutting down")
}

func (rc *RedisCache) Load() error {
	logrus.Info("rediscache loading")
	go rc.asyncPut(rc.Channel)

	cacheSize, err := rc.Redis.Get("cacheSize").Result()
	if err == redis.Nil {
		rc.Redis.Set("cacheSize", 0, 0)
	} else if err != nil {
		logrus.Errorf("Error getting cacheSize from redis: %s", err)
	}
	intCacheSize, _ := strconv.Atoi(cacheSize)
	rc.CacheSize = intCacheSize
	return nil
}

func (rc *RedisCache) Save() error {
	logrus.Info("rediscache shutting down")
	close(rc.Channel)
	rc.Redis.Close()
	return nil
}
