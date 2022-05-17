package apis_test

import (
	"bytes"
	"errors"
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/apis"
	"github.com/neda1985/getir/pkg/apis/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testInMemoryRequesInsertData = model.InMemoryDataIO{
		Key:   "active-tabs",
		Value: "getir",
	}
)

func TestInMemoryManagerService_FetchFromMemory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		testName         string
		dataManager      func() *mocks.InMemoryManager
		expectedResponse func(th *TestHelper, rec *httptest.ResponseRecorder)
		method           string
		queryString      string
	}{

		{
			testName:    "empty key",
			queryString: "",
			method:      http.MethodGet,
			dataManager: func() *mocks.InMemoryManager {
				return &mocks.InMemoryManager{}
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusBadRequest, rec.Code)
			},
		},
		{
			testName:    "InMemory fetch data return error",
			queryString: "?key=active-tabs",
			method:      http.MethodGet,
			dataManager: func() *mocks.InMemoryManager {
				m := &mocks.InMemoryManager{}
				m.On("Fetch", "active-tabs").Return("", errors.New("this is an error"))
				return m
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusInternalServerError, rec.Code)
			},
		},
		{
			testName:    "success",
			queryString: "?key=active-tabs",
			method:      http.MethodGet,
			dataManager: func() *mocks.InMemoryManager {
				m := &mocks.InMemoryManager{}
				m.On("Fetch", "active-tabs").Return("active-tabs", nil)
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
			req := httptest.NewRequest(test.method, "/get-in-memory"+test.queryString, nil)
			inMemoryManager := test.dataManager()
			m := apis.InMemoryManagerService{InMemory: inMemoryManager}
			m.FetchFromMemory(wr, req)
			test.expectedResponse(th, wr)
			inMemoryManager.AssertExpectations(t)
		})
	}

}
func TestInMemoryManagerService_PostInMemory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		testName         string
		inMemoryManager  func() *mocks.InMemoryManager
		expectedResponse func(th *TestHelper, rec *httptest.ResponseRecorder)
		method           string
		requestBody      []byte
	}{
		{
			testName:    "empty body",
			requestBody: nil,
			method:      http.MethodPost,
			inMemoryManager: func() *mocks.InMemoryManager {
				return &mocks.InMemoryManager{}
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusBadRequest, rec.Code)
			},
		},
		{
			testName: "in memory manager post data return error",
			requestBody: []byte(`{
                        "key": "active-tabs",
                        "value": "getir"}`),
			method: http.MethodPost,
			inMemoryManager: func() *mocks.InMemoryManager {
				m := &mocks.InMemoryManager{}
				m.On("Insert", testInMemoryRequesInsertData).Return(errors.New("this is an error"))
				return m
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusInternalServerError, rec.Code)
			},
		},
		{
			testName: "success",
			requestBody: []byte(`{
                        "key": "active-tabs",
                        "value": "getir"}`),
			method: http.MethodPost,
			inMemoryManager: func() *mocks.InMemoryManager {
				m := &mocks.InMemoryManager{}
				m.On("Insert", testInMemoryRequesInsertData).Return(nil)
				return m
			},
			expectedResponse: func(th *TestHelper, rec *httptest.ResponseRecorder) {
				th.Equal(http.StatusCreated, rec.Code)
			},
		},
	}
	th := New(t)
	for i := range tests {
		test := tests[i]
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			wr := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, "/set-in-memory", bytes.NewBuffer(test.requestBody))
			inMemoryManager := test.inMemoryManager()
			m := apis.InMemoryManagerService{InMemory: inMemoryManager}
			m.PostInMemory(wr, req)
			test.expectedResponse(th, wr)
			inMemoryManager.AssertExpectations(t)
		})
	}

}
