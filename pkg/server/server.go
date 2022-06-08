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

func (s *Servinator) fooHandler() gin.HandlerFunc {

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
	files, err := s.getHtmlFiles()
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

func (s *Servinator) getHtmlFiles() (*[]string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	rootDirFolders := strings.Split(workdir, string(os.PathSeparator))
	log.Debugf("rootDirsFolders: %v, %d", rootDirFolders, len(rootDirFolders))
	rootDirFolders = rootDirFolders[:len(rootDirFolders)-2]
	rootDirPath := strings.Join(rootDirFolders, "")
	log.Debugf("root path: %s", rootDirPath)

	err = os.Chdir(rootDirPath)
	if err != nil {
		return nil, err
	}

	dir, err := os.ReadDir("html")
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(dir))
	for i, entry := range dir {
		fileNames[i] = entry.Name()
	}

	return &fileNames, nil
}
