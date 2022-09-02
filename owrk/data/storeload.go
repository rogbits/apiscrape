package data

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
)

func (store *OWrkStore) LoadConfigIntoStore() error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.New("could not get working directory")
	}

	fp := path.Join(wd, "config")
	file, err := os.Open(fp)
	if err != nil {
		return errors.New("could not open config file")
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		switch {
		case strings.HasPrefix(word, "ms_between_requests="):
			value := strings.TrimLeft(word, "ms_between_requests=")
			store.MillisecondsBetweenRequests, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	err = scanner.Err()
	return err
}
