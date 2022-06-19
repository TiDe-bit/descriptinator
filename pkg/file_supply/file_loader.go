package file_supply

import (
	"github.com/sirupsen/logrus"
	"os"
)

type FileData *[]byte

func LoadFile(fileName string) (FileData, bool) {
	file, err := os.Open(fileName)
	if err != nil {
		logrus.WithError(err).Debugf("can't load file %s", fileName)
		return nil, false
	}
	defer file.Close()

	var buffer []byte

	_, err = file.Read(buffer)
	if err != nil {
		return nil, false
	}

	return &buffer, true
}

func LoadLegalText(custom ...string) *string {
	if len(custom) > 0 {
		return nil
	}
	return nil
}

func LoadAuctionText(custom ...string) *string {
	return nil
}

func LoadSellerText(custom ...string) *string {
	return nil
}

func LoadBriefText() *string {

	return nil
}

func LoadPaketText() *string {

	return nil
}

func LoadPaketBrieftaube() *string {

	return nil
}
