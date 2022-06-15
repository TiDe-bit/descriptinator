package file_supply

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func GetDescription(path string, marshaller Marshaller) (FileData, error) {
	cachedData, ok := LoadFile(path)
	if !ok {
		data, err := marshaller.CreatDescription()
		if err != nil {
			return nil, err
		}
		cachedData = data
	}

	return cachedData, nil
}

func gotoTmpl() (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	log.Debugf("WD: %s", workdir)

	rootDirFolders := strings.Split(workdir, string(os.PathSeparator))
	log.Debugf("rootDirsFolders: %v, %d", rootDirFolders, len(rootDirFolders))

	rootDirFolders = rootDirFolders[:len(rootDirFolders)-2]

	rootDirPath := strings.Join(rootDirFolders, string(os.PathSeparator))
	log.Debugf("root path: %s", rootDirPath)

	err = os.Chdir(rootDirPath)
	if err != nil {
		return "", err
	}

	return rootDirPath, nil
}

const templateFolderName = "template"

func getTmplFile(rootPath string) (string, error) {
	files, err := os.ReadDir(templateFolderName)
	if err != nil {
		return "", err
	}

	fullPath := strings.Join(
		[]string{
			rootPath,
			templateFolderName,
			files[0].Name(),
		},
		string(os.PathSeparator),
	)

	return fullPath, nil
}
