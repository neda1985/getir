package apis

import (
	"encoding/json"
	"errors"
	"github.com/neda1985/getir/model"
	"github.com/neda1985/getir/pkg/logger"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	errBadRequest         = errors.New("bad request")
	errNotFoundKey        = errors.New("key not found")
	errKeyValueRequired   = errors.New("key and value are required")
	errEmptyBody          = errors.New("request body can not be empty")
	errNotFound           = errors.New("not found")
	errKeyRequired        = errors.New("key should be set")
	errValueCanNotBeEmpty = errors.New("please fill required data")
)

//go:generate mockery --name MongoManager
type MongoManager interface {
	FetchResponse(model.FetchRequest) (response model.FetchResponse, err error)
}
type MongoMangerService struct {
	MongoManager MongoManager
}

func (d *MongoMangerService) FetchData(rw http.ResponseWriter, request *http.Request) {
	var response model.FetchResponse
	if request.Method != http.MethodPost {
		handelErrors(rw, http.StatusBadRequest, errBadRequest)
		return
	}
	if nil == request.Body {
		handelErrors(rw, http.StatusBadRequest, errEmptyBody)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.LogError(err)
		}
	}(request.Body)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handelErrors(rw, http.StatusBadRequest, err)
		return
	}

	var fetchRequest model.FetchRequest
	if err := json.Unmarshal(body, &fetchRequest); err != nil {
		handelErrors(rw, http.StatusBadRequest, err)
		return
	}
	err = validateFetchStruct(fetchRequest)
	if err != nil {
		handelErrors(rw, http.StatusBadRequest, errValueCanNotBeEmpty)
		return
	}

	response, err = d.MongoManager.FetchResponse(fetchRequest)
	if err != nil {
		handelErrors(rw, http.StatusInternalServerError, err)
		return
	}
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		logger.LogError(err)
		handelErrors(rw, http.StatusNotFound, errNotFound)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
func handelErrors(rw http.ResponseWriter, code int, message error) {
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(model.ApiError{
		Message: message.Error(),
	}); err != nil {
		logger.LogError(err)
	}

}
func validateFetchStruct(e model.FetchRequest) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}
