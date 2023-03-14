package file_supply

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Entry struct {
	ArtikelNum string  `bson:"artikelNum" json:"artikelNum"`
	Title      string  `bson:"title" json:"title"`
	Subtitle   string  `bson:"subtitle" json:"subtitle"`
	Article    Article `bson:"article" json:"article"`
	Shipping   string  `bson:"shipping" json:"shipping"`
	Legal      string  `bson:"legal" json:"legal"`
	Auction    string  `bson:"auction" json:"auction"`
	Seller     string  `bson:"seller" json:"seller"`
	Dsgvo      string  `bson:"dsgvo" json:"dsgvo"`
}

func (e *Entry) Byte() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Error(err)
		return nil
	}
	return bytes
}

func (e *Entry) WithTitle(title string) {
	e.Title = title
}
func (e *Entry) WithSubtitle(subtitle string) {
	e.Subtitle = subtitle
}

type IFileLoader interface {
}

type FileData *[]byte

// ToDo: save versandart expecially
type ITextLoader interface {
	// InitiateAction sets the article number; if not set the default gets loaded
	InitiateAction(artikelNum string) ITextLoader
	LoadTitleText(ctx context.Context) *string
	LoadLegalText(ctx context.Context) *string
	LoadAuctionText(ctx context.Context) *string
	LoadSellerText(ctx context.Context) *string

	LoadEntry(ctx context.Context, shippingMathod string) *Entry
	SaveDefaultVersandText(ctx context.Context, method Versand, data string) error
	LoadShippingText(ctx context.Context, shippingMethod string) *string
}

type Article struct {
	GeneralInfo string   `bson:"generalInfo"`
	Description string   `bson:"description"`
	Fitting     string   `bson:"fitting"`
	Condition   string   `bson:"condition"`
	Shipping    Shipping `bson:"shipping"`
}

type Shipping string

type Auction string
type Legal string
type Seller string
type Dsgvo string

type Versand string

const (
	VersandPaket      Versand = "paket"
	VersandBrief      Versand = "brief"
	VersandBrieftaube Versand = "brieftaube"
)

func (v Versand) String() string {
	return string(v)
}
