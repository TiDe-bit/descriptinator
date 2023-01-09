package file_supply

import "context"

var _ ITextLoader = &DemoLoader{}

type DemoLoader struct {
	articles map[string]DemoArticle
}

type DemoArticle struct {
	Title   string
	Legal   string
	Auction string
	Seller  string
	Brief   string
	Paket   string
	Taube   string
}

func (d *DemoLoader) Initialte(artikelNum string) {
	if d.articles == nil {
		d.articles = make(map[string]DemoArticle)
	}

	d.articles[artikelNum] = DemoArticle{
		Title:   "demo",
		Legal:   "demo",
		Auction: "demo",
		Seller:  "demo",
		Brief:   "demo",
		Paket:   "demo",
		Taube:   "demo",
	}
}

func (d *DemoLoader) LoadTitleText(ctx context.Context) *string {
	x := ""
	return &x
}

func (d *DemoLoader) LoadLegalText(ctx context.Context, custom ...string) *string {
	x := ""
	return &x
}

func (d *DemoLoader) LoadAuctionText(ctx context.Context, custom ...string) *string {
	x := ""
	return &x
}

func (d *DemoLoader) LoadSellerText(ctx context.Context, custom ...string) *string {
	x := ""
	return &x
}

func (d *DemoLoader) LoadBriefText(ctx context.Context) *string {
	x := ""
	return &x
}

func (d *DemoLoader) LoadPaketText(ctx context.Context) *string {
	x := ""
	return &x
}

func (d *DemoLoader) LoadPaketBrieftaube(ctx context.Context) *string {
	x := ""
	return &x
}
