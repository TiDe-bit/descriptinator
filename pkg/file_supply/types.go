package file_supply

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Entry struct {
	ArtikelNum string   `bson:"artikelNum" json:"artikelNum"`
	Title      *string  `bson:"title" json:"title"`
	Subtitle   *string  `bson:"subtitle" json:"subtitle"`
	Article    Article  `bson:"article" json:"article"`
	Shipping   Shipping `bson:"shipping" json:"shipping"`
	Legal      Legal    `bson:"legal" json:"legal"`
	Auction    Auction  `bson:"auction" json:"auction"`
	Seller     Seller   `bson:"seller" json:"seller"`
	Dsgvo      Dsgvo    `bson:"dsgvo" json:"dsgvo"`
}

var _ Valid = &Entry{}

func (e *Entry) Byte() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Error(err)
		return nil
	}
	return bytes
}

func (e *Entry) WithTitle(title *string) {
	e.Title = title
}
func (e *Entry) WithSubtitle(subtitle *string) {
	e.Subtitle = subtitle
}
func (e *Entry) WithShipping(ctx context.Context, shipping Versand, l ITextLoader) {
	switch shipping {
	case VersandBrief:
		e.Shipping = l.LoadBriefText(ctx)
		break
	case VersandPaket:
		e.Shipping = l.LoadPaketText(ctx)
		break
	case VersandBrieftaube:
		e.Shipping = l.LoadPaketBrieftaube(ctx)
		break
	}
}

type IFileLoader interface {
}

type FileData *[]byte

type ITextLoader interface {
	Initialte(artikelNum string)
	LoadTitleText(ctx context.Context) *string
	LoadLegalText(ctx context.Context, custom ...string) *string
	LoadAuctionText(ctx context.Context, custom ...string) *string
	LoadSellerText(ctx context.Context, custom ...string) *string
	LoadBriefText(ctx context.Context) *string
	LoadPaketText(ctx context.Context) *string
	LoadPaketBrieftaube(ctx context.Context) *string
}

type Article struct {
	GeneralInfo *string  `bson:"generalInfo"`
	Description *string  `bson:"description"`
	Fitting     *string  `bson:"fitting"`
	Condition   *string  `bson:"condition"`
	shipping    Shipping `bson:"shipping"`
}

type Shipping *string

type Auction *string
type Legal *string
type Seller *string
type Dsgvo *string

type Versand string

const (
	VersandPaket      Versand = "paket"
	VersandBrief      Versand = "brief"
	VersandBrieftaube Versand = "brieftaube"
)

func (v Versand) String() string {
	return string(v)
}
