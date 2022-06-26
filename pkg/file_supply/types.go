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
