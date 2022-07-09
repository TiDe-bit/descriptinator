package file_supply

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

var _ ITextLoader = &MongoTextLoader{}

func ConnectToMongodb(ctx context.Context) (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongo:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoTextLoader(ctx context.Context) ITextLoader {
	client, err := ConnectToMongodb(ctx)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	db := client.Database("texts")

	err = db.CreateCollection(ctx, "defaults", options.CreateCollection())
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	defaultsCollection := db.Collection("defaults")

	err = db.CreateCollection(ctx, "articles", options.CreateCollection())
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	articlesCollection := db.Collection("articles")

	return &MongoTextLoader{
		defaultsCollection,
		articlesCollection,
		"",
	}
}

type MongoTextLoader struct {
	defaultsCollection *mongo.Collection
	articleCollection  *mongo.Collection
	currentArtikelNum  string
}

func (l *MongoTextLoader) Initialte(artikelNum string) {
	l.currentArtikelNum = artikelNum
}

func (l *MongoTextLoader) LoadTitleText(ctx context.Context) *string {
	var target *marshaller.Entry
	filter := bson.M{"artikelNum": l.currentArtikelNum}

	err := l.articleCollection.FindOne(ctx, filter).Decode(target)
	if err != nil {
		return nil
	}

	return target.Title
}

func (l *MongoTextLoader) LoadLegalText(ctx context.Context, custom ...string) *string {
	if len(custom) > 0 {
		return nil
	}

	var target marshaller.Legal
	filter := bson.M{"": ""}

	err := l.defaultsCollection.FindOne(ctx, filter).Decode(target)
	if err != nil {
		return nil
	}

	return target
}

func (l *MongoTextLoader) LoadAuctionText(ctx context.Context, custom ...string) *string {
	if len(custom) > 0 {
		return nil
	}

	var target marshaller.Auction
	filter := bson.M{"": ""}

	err := l.defaultsCollection.FindOne(ctx, filter).Decode(target)
	if err != nil {
		return nil
	}

	return target
}

type Ttext struct {
	text string `bson:"Ttext"`
}

func (l *MongoTextLoader) LoadSellerText(ctx context.Context, custom ...string) *string {
	if len(custom) > 0 {
		customText := strings.Join(custom, "")
		return &customText
	}

	filter := bson.M{"": ""}

	result := l.defaultsCollection.FindOne(ctx, filter)

	var text Ttext
	err := result.Decode(&text)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return &text.text
}

func (l *MongoTextLoader) LoadBriefText(ctx context.Context) *string {
	return nil
}

func (l *MongoTextLoader) LoadPaketText(ctx context.Context) *string {
	return nil
}

func (l *MongoTextLoader) LoadPaketBrieftaube(ctx context.Context) *string {
	return nil
}
