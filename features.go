package dp_golearn_to_rank

import (
	"net/http"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"errors"
)

const (
	LearnToRankApi = "/_ltr"
)

type FeatureService struct {
	client *Client
}

func NewFeatureService(client *Client) (*FeatureService) {
	f := FeatureService{
		client: client,
	}
	return &f
}

func (service *FeatureService) performSimpleRequest(method, path string, ignoreErrors ...int) (*elastic.Response, error) {
	return service.client.PerformRequest(method, path, nil, nil, ignoreErrors...)
}

func (service *FeatureService) DefaultFeatureStoreExists() (bool, error) {
	return service.FeatureStoreExists("")
}

func (service *FeatureService) FeatureStoreExists(featureStore string) (bool, error) {
	// Checks if the ltr feature store exists
	method := http.MethodGet
	path := LearnToRankApi
	if featureStore != "" {
		path = path + "/" + featureStore
	}

	// Perform the request - note that a 404 indicates the feature store does not exist
	// and is a perfectly valid response code
	resp, err := service.performSimpleRequest(method, path, http.StatusNotFound)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		err := errors.New(fmt.Sprintf("error checking if feature store %s exists", path))
		return false, err
	}

	exists := resp.StatusCode == http.StatusOK;

	return exists, err
}

func (service *FeatureService) CreateFeatureStore(featureStore string) (*elastic.Response, error) {
	// Creates a feature store
	method := http.MethodPut
	path := LearnToRankApi + "/" + featureStore

	resp, err := service.performSimpleRequest(method, path)

	return resp, err
}

func (service *FeatureService) DropFeatureStore(featureStore string) (*elastic.Response, error) {
	method := http.MethodDelete
	path := LearnToRankApi + "/" + featureStore

	resp, err := service.performSimpleRequest(method, path)

	return resp, err
}
