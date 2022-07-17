package api

import (
	"descriptinator/pkg/file_supply"
	"descriptinator/pkg/marshaller"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Run(engine *gin.Engine) {
	// ToDo api path
	engine.GET(marshaller.VERSAND_BRIEF.String(), loadVersand(marshaller.VERSAND_BRIEF.String()))
	engine.GET(marshaller.VERSAND_PAKET.String(), loadVersand(marshaller.VERSAND_PAKET.String()))
	engine.GET(marshaller.VERSAND_BRIEFTAUBE.String(), loadVersand(marshaller.VERSAND_BRIEFTAUBE.String()))

	// ToDo do other imports
	engine.PUT(marshaller.VERSAND_BRIEF.String(), saveVersand(marshaller.VERSAND_BRIEF.String()))
	engine.PUT(marshaller.VERSAND_PAKET.String(), saveVersand(marshaller.VERSAND_PAKET.String()))
	engine.PUT(marshaller.VERSAND_BRIEFTAUBE.String(), saveVersand(marshaller.VERSAND_BRIEFTAUBE.String()))
}

func loadVersand(versandArt string) gin.HandlerFunc {
	return load[file_supply.Ttext]("versand", versandArt, false)
}

func saveVersand(versandArt string) gin.HandlerFunc {
	return save[file_supply.Ttext]("versand", versandArt, false)
}

func loadEntry(entryID string) gin.HandlerFunc {
	return load[*marshaller.Entry]("entry", entryID, true)
}

func saveEntry(entryID string) gin.HandlerFunc {
	return save[*marshaller.Entry]("entry", entryID, true)
}

func load[T file_supply.Valid](field, id string, special bool) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var data *T
		var err error
		if special {
			data, err = file_supply.LoadAny[T](
				gtx,
				bson.M{field: id},
				struct{}{},
			)
		} else {
			data, err = file_supply.LoadAny[T](
				gtx,
				bson.M{field: id},
			)
		}
		if err != nil {
			gtx.Error(err)
		}
		tmp := *data
		gtx.Data(http.StatusOK, "data", tmp.Byte())
	}
}

func save[T file_supply.Valid](field, id string, special bool) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var object *T

		rawData, err := gtx.GetRawData()
		if err != nil {
			gtx.Error(err)
		}

		err = json.Unmarshal(rawData, object)
		if err != nil {
			gtx.Error(err)
		}

		if special {
			err = file_supply.SaveAny[T](
				gtx,
				bson.M{field: id},
				object,
				struct{}{},
			)
		} else {
			err = file_supply.SaveAny[T](
				gtx,
				bson.M{field: id},
				object,
			)
		}
		if err != nil {
			gtx.Error(err)
		}

		gtx.Status(http.StatusOK)
	}
}
