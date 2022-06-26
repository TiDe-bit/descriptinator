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
		client,
		defaultsCollection,
		articlesCollection,
	}
}

type MongoTextLoader struct {
	*mongo.Client
	defaultsCollection *mongo.Collection
	articleCollection  *mongo.Collection
}

func (l *MongoTextLoader) LoadLegalText(ctx context.Context, custom ...string) *string {
	if len(custom) > 0 {
		return nil
	}
	return nil
}

func (l *MongoTextLoader) LoadAuctionText(ctx context.Context, custom ...string) *string {
	if len(custom) > 0 {
		return nil
	}
	return nil
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
