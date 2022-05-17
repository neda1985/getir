package apis_test

import (
	"bytes"
	"errors"
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/apis"
	"github.com/neda1985/getir/pkg/apis/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testResponseModel = model.FetchResponse{
		Code:    1,
		Msg:     "test",
		Records: []model.Records{},
	}
	testRequestFetchData = model.FetchRequest{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}
)

func TestMongoMangerService_FetchData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		testName         string
		mongoDataManager func() *mocks.MongoManager
		expectedResponse func(th *TestHelper, rec *httptest.ResponseRecorder)
		method           string
		requestBody      []byte
	}{
		{
			testName:    "invalid method",
			method:      http.MethodGet,
			requestBody: nil,
			mongoDataManager: func() *mocks.MongoManager {
				return &mocks.MongoManager{}
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusBadRequest, rec.Code)
			},
		},
		{
			testName:    "empty body",
			requestBody: nil,
			method:      http.MethodPost,
			mongoDataManager: func() *mocks.MongoManager {
				return &mocks.MongoManager{}
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusBadRequest, rec.Code)
			},
		},
		{
			testName: "mongo manager fetch data return error",
			requestBody: []byte(`{
                         "startDate": "2016-01-26",
                         "endDate": "2018-02-02",
                         "minCount": 2700,
                         "maxCount": 3000}`),
			method: http.MethodPost,
			mongoDataManager: func() *mocks.MongoManager {
				m := &mocks.MongoManager{}
				m.On("FetchResponse", testRequestFetchData).Return(model.FetchResponse{}, errors.New("this is an error"))
				return m
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusInternalServerError, rec.Code)
			},
		},
		{
			testName: "success",
			requestBody: []byte(`{
                         "startDate": "2016-01-26",
                         "endDate": "2018-02-02",
                         "minCount": 2700,
                         "maxCount": 3000}`),
			method: http.MethodPost,
			mongoDataManager: func() *mocks.MongoManager {
				m := &mocks.MongoManager{}
				m.On("FetchResponse", testRequestFetchData).Return(testResponseModel, nil)
				return m
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusOK, rec.Code)
			},
		},
	}
	th := New(t)
	for i := range tests {
		test := tests[i]
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			wr := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, "/fetch", bytes.NewBuffer(test.requestBody))
			mongoDataManager := test.mongoDataManager()
			m := apis.MongoMangerService{MongoManager: mongoDataManager}
			m.FetchData(wr, req)
			test.expectedResponse(th, wr)
			mongoDataManager.AssertExpectations(t)
		})
	}

}
func New(t *testing.T) *TestHelper {
	th := TestHelper{
		Test:       t,
		Assertions: require.New(t),
	}
	return &th
}

type TestHelper struct {
	Test *testing.T
	*require.Assertions
}
