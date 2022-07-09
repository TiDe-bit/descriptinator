package file_supply

import "context"

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
