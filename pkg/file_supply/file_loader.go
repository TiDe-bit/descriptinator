package file_supply

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	defaultsCollection, articlesCollection, err := setupDB(ctx, client)
	if err != nil {
		logrus.Fatal(err)
	}

	return &MongoTextLoader{
		defaultsCollection,
		articlesCollection,
		"",
	}
}

func setupDB(ctx context.Context, client *mongo.Client) (defaultsCollection *mongo.Collection, articlesCollection *mongo.Collection, err error) {
	db := client.Database("texts")

	err = db.CreateCollection(ctx, "defaults", options.CreateCollection())
	if err != nil {
		return nil, nil, err
	}
	defaultsCollection = db.Collection("defaults")

	err = db.CreateCollection(ctx, "articles", options.CreateCollection())
	if err != nil {
		return nil, nil, err
	}
	articlesCollection = db.Collection("articles")

	return defaultsCollection, articlesCollection, nil
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
	var target *Entry
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

	var target Legal
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

	var target Auction
	filter := bson.M{"": ""}

	err := l.defaultsCollection.FindOne(ctx, filter).Decode(target)
	if err != nil {
		return nil
	}

	return target
}

type Valid interface {
	Byte() []byte
}

var _ Valid = Ttext{}

type Ttext struct {
	text string `bson:"Ttext"`
}

func (t Ttext) Byte() []byte {

	bytes, err := json.Marshal(t)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return bytes
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

func LoadAny[T Ttext | any](ctx context.Context, filter bson.M, specific ...struct{}) (*T, error) {
	var result T

	client, err := ConnectToMongodb(ctx)
	defaultsCollection, articleCollection, err := setupDB(ctx, client)
	if err != nil {
		return nil, err
	}

	if len(specific) == 0 {
		if err := defaultsCollection.FindOne(ctx, filter).Decode(&result); err != nil {
			return nil, err
		}
	} else {
		if err := articleCollection.FindOne(ctx, filter).Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func SaveAny[T any](ctx context.Context, filter bson.M, data *T, specific ...struct{}) error {
	// ToDo: only one connection...
	client, err := ConnectToMongodb(ctx)
	defaultsCollection, articleCollection, err := setupDB(ctx, client)
	opts := options.Update().SetUpsert(true)

	if err != nil {
		return err
	}

	if len(specific) == 0 {
		_, err := defaultsCollection.UpdateOne(ctx, filter, data, opts)
		if err != nil {
			return err
		}
	} else {
		_, err := articleCollection.UpdateOne(ctx, filter, data, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
