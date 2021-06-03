package mongodb

import (
	"context"
	"time"

	"github.com/inblack67/url-shortner/shortener"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongodbRepo struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

var (
	collectionName = "redirects"
)

func newMongodbClient(mongoURL string, mongodbTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongodbTimeout))
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	pingErr := client.Ping(ctx, readpref.Primary()) // confirms that we can read our db or not
	if pingErr != nil {
		return nil, pingErr
	}
	return client, nil
}

func NewMongodbRepo(mongodbURL, database string, mongodbTimeout int) (shortener.RedirectRepo, error) {
	repo := &mongodbRepo{
		timeout:  time.Duration(mongodbTimeout) * time.Second,
		database: database,
	}
	client, err := newMongodbClient(mongodbURL, mongodbTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repo.NewMongodbRepo")
	}
	repo.client = client
	return repo, nil
}

func (r *mongodbRepo) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	redirect := &shortener.Redirect{}

	collection := r.client.Database(r.database).Collection(collectionName)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrorRedirectNotFound, "repo.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repo.Redirect.Find")
	}

	return redirect, nil
}

func (r *mongodbRepo) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	collection := r.client.Database(r.database).Collection(collectionName)

	_, err := collection.InsertOne(ctx, bson.M{
		"code": redirect.Code,
		"ur;": redirect.URL,
		"createdAt": redirect.CreatedAt
	})

}
