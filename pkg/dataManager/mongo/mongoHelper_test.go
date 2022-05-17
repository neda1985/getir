package mongodb_test

import (
	"github.com/neda1985/getir/model"
	mongodb "github.com/neda1985/getir/pkg/dataManager/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var (
	testMongoQueryInput = model.FetchRequest{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}
	t1       = time.Date(2018, time.February, 2, 0, 0, 0, 0, time.UTC)
	t2       = time.Date(2016, time.January, 26, 0, 0, 0, 0, time.UTC)
	expected = []primitive.M([]primitive.M{{"$match": primitive.M{
		"createdAt": bson.M{
			"$gt": &t2,
			"$lt": &t1,
		},
	},
	},
		{"$project": primitive.M{"_id": 0, "createdAt": 1, "key": 1, "totalCount": primitive.M{"$sum": "$counts"}}}, {"$match": mongodb.TotalCount{TotalCount: mongodb.GTQuery{Gt: float64(2700), Lt: float64(3000)}}}})
)

type CreatedAt struct {
	createdAt mongodb.GTQuery `bson:"createdAt,omitempty"`
}

func TestMongoHelperService_BuildFetchQuery(t *testing.T) {
	query, err := mongodb.BuildFetchQuery(testMongoQueryInput)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, expected, query)
}
