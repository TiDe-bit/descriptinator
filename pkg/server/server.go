package server

import (
	"descriptinator/pkg/file_supply"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Servinator struct {
	address string
}

func NewServinator(address string) *Servinator {
	return &Servinator{
		address: address,
	}
}

func (s *Servinator) Serve() {
	engine := gin.Default()

	engine.Group("brief", handleGroupBrief(engine))

	err := engine.Run(s.address)
	if err != nil {
		log.Fatal(err)
	}
}

func handleGroupBrief(engine *gin.Engine) func(g *gin.Context) {
	engine.GET("")
}

func sendDescription(data file_supply.FileData) func(g *gin.Context) {
	return func(g *gin.Context) {
		g.Data(
			http.StatusOK,
			http.DetectContentType(*data),
			*data,
		)
	}

}
