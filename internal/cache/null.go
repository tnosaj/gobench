package cache

import (
	"github.com/samborkent/uuid"
	"github.com/sirupsen/logrus"
)

type NullCache struct {
}

func newNullCache() *NullCache {
	logrus.Info("Starting with no cache")
	return &NullCache{}
}

func (fc *NullCache) GetRandom() (uuid.UUID, error) {
	logrus.Trace("nullcache get random")
	return uuid.New(), nil
}

func (fc *NullCache) Put(uuid uuid.UUID) error {
	logrus.Tracef("nullcache random: %s", uuid)
	return nil
}

func (nc *NullCache) Load() error {
	logrus.Debug("nullcache not loading anything")
	return nil
}

func (nc *NullCache) Save() error {
	logrus.Debug("nullcache not saving anything")
	return nil
}
