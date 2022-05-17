package mongodb

import (
	"context"
	"fmt"
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

const (
	ConnectionString = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"
	DateLayout       = "2006-01-02"
	DB               = "getir-case-study"
	Collection       = "records"
	NoDataFound      = "empty record"
	Success          = "Success"
)

func NewMongoInstance() DataManager {
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(ConnectionString))
	if err != nil {
		logger.LogError(err)
	}
	session, err := client.StartSession()
	if err != nil {
		logger.LogError(err)
	}
	database := session.Client().Database(DB)

	recordsCollection := database.Collection(Collection)
	return &mongodb{collection: recordsCollection}
}

type DataManager interface {
	FetchResponse(model.FetchRequest) (response model.FetchResponse, err error)
}
type mongodb struct {
	collection *mongo.Collection
}

func (m *mongodb) FetchResponse(input model.FetchRequest) (response model.FetchResponse, err error) {
	var fetchResponse model.FetchResponse
	fetchRequest := input
	fetchResponse.Records = []model.Records{}
	query, err := BuildFetchQuery(fetchRequest)
	if err != nil {
		fetchResponse.Msg = err.Error()
		logger.LogError(err)
		return fetchResponse, err
	}
	cursor, err := m.collection.Aggregate(context.TODO(), query)
	if err != nil {
		fetchResponse.Msg = err.Error()
		logger.LogError(err)
		return fetchResponse, err
	}

	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &fetchResponse.Records); err != nil {
		fetchResponse.Msg = err.Error()
		logger.LogError(err)
		return fetchResponse, err
	}

	if len(fetchResponse.Records) > 0 {
		fetchResponse.Code = 0
		fetchResponse.Msg = Success
		return fetchResponse, nil
	}

	fetchResponse.Code = http.StatusNoContent
	fetchResponse.Msg = NoDataFound
	err = fmt.Errorf(NoDataFound)
	return fetchResponse, err
}
