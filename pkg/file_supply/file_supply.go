package file_supply

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func FilePathFromArtikelNr(artikelNr string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	rootDirFolders := strings.Split(wd, string(os.PathSeparator))

	rootDirFolders = rootDirFolders[:len(rootDirFolders)-2]

	rootDirPath := strings.Join(rootDirFolders, string(os.PathSeparator))

	fileName := strings.Join([]string{artikelNr, "html"}, ".")
	filePath := strings.Join([]string{rootDirPath, "html", fileName}, "/")
	return filePath
}

func GotoTmpl() (string, error) {
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

func GetTmplFile(rootPath string) (string, error) {
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

// ToDo: buffer in map
func LoadFile(fileName string) (FileData, bool) {
	file, err := os.Open(fileName)
	if err != nil {
		log.WithError(err).Debugf("can't load file %s", fileName)
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
