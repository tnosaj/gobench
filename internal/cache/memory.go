package cache

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/samborkent/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

type MemoryCache struct {
	filePath   string
	Randomizer internal.Random
	Channel    chan uuid.UUID
	Cache      []string
}

func newMemoryCache(path string, randomizer internal.Random) *MemoryCache {
	logrus.Infof("Starting new file cache for %s", path)
	channel := make(chan uuid.UUID)
	return &MemoryCache{filePath: path, Randomizer: randomizer, Channel: channel}
}

func (fc *MemoryCache) GetRandom() (uuid.UUID, error) {
	logrus.Trace("MemoryCache get random")
	randomIndex := fc.Randomizer.Intn(len(fc.Cache))
	uid, _ := uuid.StringToUUID(fc.Cache[randomIndex])
	return uid, nil
}

func (fc *MemoryCache) Put(uuid uuid.UUID) error {
	logrus.Tracef("MemoryCache put: %s", uuid)
	fc.Channel <- uuid
	return nil
}

func (fc *MemoryCache) asyncPut(c chan uuid.UUID) {

	for uuid := range c {
		fc.Cache = append(fc.Cache, uuid.String())
	}
}

func (fc *MemoryCache) Save() error {
	close(fc.Channel)
	if fc.filePath != "none" {
		err := writeLines(fc.Cache, fc.filePath)
		if err != nil {
			return fmt.Errorf("Could not write cache to disk because %s", err)
		}
		logrus.Info("File Closed, shutting down")
		return nil
	}
	logrus.Info("No file defined, shutting down")
	return nil
}

func (fc *MemoryCache) Load() error {
	logrus.Infof("Start loading MemoryCache")
	go fc.asyncPut(fc.Channel)
	var values []string

	if fc.filePath != "none" {
		if _, err := os.Stat(fc.filePath); errors.Is(err, os.ErrNotExist) {
			logrus.Debug("Tmpfile does not exist on startup")
		} else {
			logrus.Debugf("Tmpfile exists on startup, reading")
			values, err = readLines(fc.filePath)
			if err != nil {
				logrus.Errorf("There was a problem reading the tmpfile into memory: %s", err)
			}
		}
	} else {
		logrus.Info("No file set - no cache persistence between restarts")
	}
	fc.Cache = values

	logrus.Infof("Finished loading %d items into MemoryCache", len(fc.Cache))
	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	var file *os.File
	var err error

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logrus.Debug("Tmpfile does not exist on shutdown, creating")
		file, err = os.Create(path)
	} else {
		logrus.Debugf("Tmpfile exists on shutdown, appending")
		file, err = os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0600)
	}

	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
