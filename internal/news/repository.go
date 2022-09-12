package news

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go-news-feed/pkg/model"
)

const (
	orderDesc = "desc"
	maxLimit  = 1000
)

// Repository - interface
//
//go:generate mockgen -source=repository.go -destination=repository_mock.go --package=news
type Repository interface {
	FindByID(ctx context.Context, id string) (model.Article, error)
	Find(ctx context.Context, fr model.FindRequest) (model.FindResponse, error)
	Create(ctx context.Context, article model.Article) error
}

type repository struct {
	collection *mongo.Collection
}

// newRepository - constructor
func newRepository(ctx context.Context, config MongoConfig) (Repository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	collection := client.Database(config.Database).Collection(config.Collection)

	return &repository{collection: collection}, nil
}

func (r repository) FindByID(ctx context.Context, id string) (model.Article, error) {
	var article model.Article

	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article); err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (r repository) Find(ctx context.Context, fr model.FindRequest) (model.FindResponse, error) {
	pipeline := mongo.Pipeline{}

	if fr.Category != "" {
		pipeline = append(pipeline, r.buildFilterStage("source.category", fr.Category))
	}

	if fr.Provider != "" {
		pipeline = append(pipeline, r.buildFilterStage("source.provider", fr.Provider))
	}

	if fr.Sort != "" {
		pipeline = append(pipeline, r.buildOrderStage(fr.Sort, fr.Order))
	}

	pipeline = append(pipeline,
		r.buildFacetStage(fr.Page, fr.Limit),
		r.buildProjectStage(),
	)

	response, err := r.aggregate(ctx, pipeline)
	if err != nil {
		return model.FindResponse{}, err
	}

	response.Criteria = fr

	return response, nil
}

func (r repository) Create(ctx context.Context, article model.Article) error {
	_, err := r.collection.InsertOne(ctx, &article)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) aggregate(ctx context.Context, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) (model.FindResponse, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return model.FindResponse{}, err
	}

	response := make([]model.FindResponse, 0)
	if err := cursor.All(ctx, &response); err != nil {
		return model.FindResponse{}, err
	}

	return response[0], nil
}

// buildFilterStage used to filter by field provided
func (r repository) buildFilterStage(field, value string) bson.D {
	return bson.D{
		{
			Key: "$match", Value: bson.D{
				{Key: field, Value: value},
			},
		},
	}
}

// buildOrderStage used to process a order stage as part of the
// aggregation pipeline
func (r repository) buildOrderStage(sort, order string) bson.D {
	orderN := 1

	if strings.ToLower(order) == orderDesc {
		orderN = -1
	}

	return bson.D{
		{
			Key: "$sort", Value: bson.D{
				{Key: sort, Value: orderN},
			},
		},
	}
}

// buildFacetStage used to process multiple aggregation pipelines within a single stage
// specifying the sub-pipeline output.
// - Count Stage
// - Pagination stage
// - Limit stage
func (r repository) buildFacetStage(page, limit int) bson.D {
	if limit == 0 || limit > maxLimit {
		limit = maxLimit
	}

	return bson.D{
		{
			Key: "$facet", Value: bson.D{
				{
					Key: "metadata", Value: bson.A{
						bson.D{
							{
								Key: "$count", Value: "total",
							},
						},
					},
				},
				{
					Key: "articles", Value: bson.A{
						bson.D{
							{
								Key: "$skip", Value: page * limit,
							},
						},
						bson.D{
							{
								Key: "$limit", Value: limit,
							},
						},
					},
				},
			},
		},
	}
}

// buildProjectStage customise the output
func (r repository) buildProjectStage() bson.D {
	return bson.D{
		{
			Key: "$project", Value: bson.D{
				{
					Key: "total", Value: bson.D{
						{
							Key: "$arrayElemAt", Value: bson.A{"$metadata.total", 0},
						},
					},
				},
				{
					Key: "articles", Value: 1,
				},
			},
		},
	}
}
