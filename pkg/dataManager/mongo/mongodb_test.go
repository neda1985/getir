package mongodb_test

import (
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/dataManager/mongo"
	"testing"
)

var (
	mongoManager       = mongodb.NewMongoInstance()
	testFetchDataInput = model.FetchRequest{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MaxCount:  3000,
		MinCount:  2700,
	}
)

func TestMongodb_FetchResponse(t *testing.T) {
	res, err := mongoManager.FetchResponse(testFetchDataInput)
	if err != nil {
		t.Fail()
	}
	if res.Code != 0 {
		t.Fail()
	}
}
