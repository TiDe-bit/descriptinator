package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type Servinator struct {
	gin     *gin.Engine
	address string
}

func NewServinator(address string) *Servinator {
	engine := gin.Default()
	return &Servinator{
		gin:     engine,
		address: address,
	}
}

func (s *Servinator) routing() {
	s.gin.Group("foo", func(ctx *gin.Context) {
		ctx.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})
}

func (s *Servinator) Serve() {
	files, err := GetHtmlFiles()
	if err != nil {
		log.Fatal(err)
	}

	s.gin.LoadHTMLFiles(*files...)
	s.routing()

	err = s.gin.Run(s.address)
	// ToDo: errChan or something
	if err != nil {
		log.Fatal(err)
	}
}

const htmlFolderName = "html"

func GetHtmlFiles() (*[]string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	log.Debugf("WD: %s", workdir)

	rootDirFolders := strings.Split(workdir, string(os.PathSeparator))
	log.Debugf("rootDirsFolders: %v, %d", rootDirFolders, len(rootDirFolders))

	rootDirFolders = rootDirFolders[:len(rootDirFolders)-2]

	rootDirPath := strings.Join(rootDirFolders, string(os.PathSeparator))
	log.Debugf("root path: %s", rootDirPath)

	err = os.Chdir(rootDirPath)
	if err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(htmlFolderName)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(dir))
	for i, entry := range dir {
		fileNames[i] = strings.Join(
			[]string{
				rootDirPath,
				htmlFolderName,
				entry.Name(),
			},
			string(os.PathSeparator),
		)
	}

	log.Debugf("filenames: %v", fileNames)

	return &fileNames, nil
}
