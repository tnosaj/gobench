package cache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/samborkent/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
)

type FileCache struct {
	filePath   string
	Randomizer internal.Random
	Channel    chan uuid.UUID
	CacheSize  int
}

func newFileCache(path string, randomizer internal.Random) *FileCache {
	logrus.Infof("Starting new file cache for %s", path)
	channel := make(chan uuid.UUID)
	return &FileCache{filePath: path, Randomizer: randomizer, Channel: channel}
}

func (fc *FileCache) GetRandom() (uuid.UUID, error) {
	logrus.Trace("filecache get random")
	cmd := exec.Command("sed",
		"-n", strconv.Itoa(fc.Randomizer.Intn(fc.CacheSize))+"p",
		fc.filePath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return uuid.NewV7(), fmt.Errorf("Error running 'sed' command: %s", err)
	}

	id, err := uuid.StringToUUID(out.String())
	if err != nil {
		return uuid.NewV7(), fmt.Errorf("Error converting key: %s", err)
	}
	return id, nil
}

func (fc *FileCache) Put(uuid uuid.UUID) error {
	logrus.Tracef("filecache put: %s", uuid)
	fc.Channel <- uuid
	return nil
}

func (fc *FileCache) asyncPut(c chan uuid.UUID) {
	var file *os.File
	if _, err := os.Stat(fc.filePath); errors.Is(err, os.ErrNotExist) {
		logrus.Info("Tmpfile does not exist on startup")
		file, err = os.Create(fc.filePath)
		if err != nil {
			logrus.Errorf("There was a problem opening the tmpfile: %s", err)
		}
	} else {
		logrus.Info("Tmpfile exists on startup, opening")
		file, err = os.OpenFile(fc.filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logrus.Errorf("There was a problem opening the tmpfile: %s", err)
		}
		logrus.Info("finished opening tmpfile")
	}
	w := bufio.NewWriter(file)

	for uuid := range c {
		fc.CacheSize++
		w.WriteString(uuid.String() + "\n")
	}
	w.Flush()
	file.Close()
	logrus.Info("Channel Closed, shutting down")
}

func (fc *FileCache) Load() error {
	logrus.Infof("Start loading FileCache")
	go fc.asyncPut(fc.Channel)
	fc.CacheSize = 0

	i, err := countLinesWithWC(fc.filePath)
	if err != nil {
		return fmt.Errorf("Error setting up current file length: %s", err)
	}
	logrus.Infof("Finished loading FileCache")
	fc.CacheSize = i
	return nil
}

func (fc *FileCache) Save() error {
	close(fc.Channel)

	logrus.Info("filecache: no data to write")
	return nil
}

// chatgpt
func countLinesWithWC(filename string) (int, error) {
	cmd := exec.Command("wc", "-l", filename)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	// Parse the output to get the line count
	output := strings.Fields(out.String())
	if len(output) == 0 {
		return 0, fmt.Errorf("unexpected output from wc")
	}

	lineCount, err := strconv.Atoi(output[0])
	if err != nil {
		return 0, err
	}

	return lineCount, nil
}
