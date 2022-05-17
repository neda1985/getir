package mongodb

import (
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Project struct {
	ID         int    `bson:"_id,omitempty"`
	Key        int    `bson:"key,omitempty"`
	createdAt  int    `bson:"createdAt,omitempty"`
	totalCount bson.M `bson:"totalCount,omitempty"`
}
type GTQuery struct {
	Gt interface{} `bson:"$gt,omitempty"`
	Lt interface{} `bson:"$lt,omitempty"`
}
type CreatedAt struct {
	createdAt GTQuery `bson:"createdAt,omitempty"`
}
type TotalCount struct {
	TotalCount GTQuery `bson:"totalCount,omitempty"`
}

func BuildFetchQuery(fetchRequest model.FetchRequest) ([]bson.M, error) {
	startDate, endDate, err := dateParser(fetchRequest)
	if err != nil {
		return nil, err
	}
	totalCount := GTQuery{
		Gt: fetchRequest.MinCount,
		Lt: fetchRequest.MaxCount,
	}
	query := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gt": startDate,
					"$lt": endDate,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"key":        1,
				"createdAt":  1,
				"totalCount": bson.M{"$sum": "$counts"},
			},
		},
		{
			"$match": TotalCount{GTQuery{
				Gt: totalCount.Gt,
				Lt: totalCount.Lt,
			}},
		},
	}
	return query, nil
}

func dateParser(fetchRequest model.FetchRequest) (*time.Time, *time.Time, error) {
	startDate, err := time.Parse(DateLayout, fetchRequest.StartDate)
	if err != nil {
		logger.LogError(err)
		return nil, nil, err
	}
	endDate, err := time.Parse(DateLayout, fetchRequest.EndDate)
	if err != nil {
		logger.LogError(err)
		return nil, nil, err
	}
	return &startDate, &endDate, nil
}
