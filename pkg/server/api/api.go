package api

import (
	"descriptinator/pkg/file_supply"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run(engine *gin.Engine, loader file_supply.ITextLoader) {
	engine.GET(getApiPath("entry"), func(gtx *gin.Context) {
		artikelNum, exists := gtx.Params.Get("artikelNum")
		if !exists {
			gtx.Error(errors.New("no artikelNum in path"))
			return
		}
		loadEntryFromDB(artikelNum, loader)
	})

	engine.PUT(getApiPath("entry"), func(gtx *gin.Context) {
		artikelNum, exists := gtx.Params.Get("artikelNum")
		if !exists {
			gtx.Error(errors.New("no artikelNum in path"))
			return
		}
		saveEntry(artikelNum, loader)
	})
}

func getApiPath(endPoint string) string {
	return fmt.Sprintf("edit/%s", endPoint)
}

func loadEntryFromDB(entryID string, loader file_supply.ITextLoader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loader.LoadEntry(ctx, entryID)
	}
}

func saveEntry(entryID string, loader file_supply.ITextLoader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loader.SaveEntry(ctx, nil) // ToDo
	}
}
