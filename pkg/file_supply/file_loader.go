package file_supply

import (
	"github.com/sirupsen/logrus"
	"os"
)

type FileData *[]byte

func LoadFile(fileName string) (data FileData, ok bool) {
	file, err := os.Open(fileName)
	if err != nil {
		logrus.WithError(err).Debugf("can't load file %s", fileName)
		return nil, false
	}
	defer file.Close()

	_, err = file.Read(*data)
	if err != nil {
		return nil, false
	}

	return data, true
}
