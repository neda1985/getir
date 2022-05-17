package apis

import (
	"encoding/json"
	"github.com/neda1985/getir/model"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type InMemoryManager interface {
	Insert(input model.InMemoryDataIO) error
	Fetch(key string) (string, error)
}
type InMemoryManagerService struct {
	InMemory InMemoryManager
}

func (i *InMemoryManagerService) SwitchRout(rw http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		i.FetchFromMemory(rw, request)
	case http.MethodPost:
		i.PostInMemory(rw, request)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}
func (i *InMemoryManagerService) FetchFromMemory(rw http.ResponseWriter, request *http.Request) {

	if request.URL.Query().Get("key") == "" {
		handelErrors(rw, http.StatusBadRequest, errKeyRequired)
		return
	}
	value, err := i.InMemory.Fetch(request.URL.Query().Get("key"))
	if err != nil {
		handelErrors(rw, http.StatusInternalServerError, err)
		return
	}
	if value == "" {
		handelErrors(rw, http.StatusNotFound, errNotFoundKey)
		return
	}
	result := model.InMemoryDataIO{Key: request.URL.Query().Get("key"), Value: value}
	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		handelErrors(rw, http.StatusInternalServerError, err)
		return
	}

}
func (i *InMemoryManagerService) PostInMemory(rw http.ResponseWriter, request *http.Request) {

	var input model.InMemoryDataIO
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		handelErrors(rw, http.StatusBadRequest, errBadRequest)
		return
	}
	err = validateStruct(input)
	if err != nil {
		handelErrors(rw, http.StatusBadRequest, errKeyValueRequired)
		return
	}
	err = i.InMemory.Insert(input)
	if err != nil {
		handelErrors(rw, http.StatusInternalServerError, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(rw).Encode(&input)
	if err != nil {
		handelErrors(rw, http.StatusBadRequest, errKeyValueRequired)
		return
	}
}

func validateStruct(e model.InMemoryDataIO) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}
