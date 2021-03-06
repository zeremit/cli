package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
)

type UserProvidedServiceInstanceRepository interface {
	Create(name, drainUrl string, params map[string]interface{}) (apiErr error)
	Update(serviceInstanceFields models.ServiceInstanceFields) (apiErr error)
}

type CCUserProvidedServiceInstanceRepository struct {
	config  core_config.Reader
	gateway net.Gateway
}

func NewCCUserProvidedServiceInstanceRepository(config core_config.Reader, gateway net.Gateway) (repo CCUserProvidedServiceInstanceRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CCUserProvidedServiceInstanceRepository) Create(name, drainUrl string, params map[string]interface{}) (apiErr error) {
	path := "/v2/user_provided_service_instances"

	type RequestBody struct {
		Name           string                 `json:"name"`
		Credentials    map[string]interface{} `json:"credentials"`
		SpaceGuid      string                 `json:"space_guid"`
		SysLogDrainUrl string                 `json:"syslog_drain_url"`
	}

	jsonBytes, err := json.Marshal(RequestBody{
		Name:           name,
		Credentials:    params,
		SpaceGuid:      repo.config.SpaceFields().Guid,
		SysLogDrainUrl: drainUrl,
	})

	if err != nil {
		apiErr = errors.NewWithError("Error parsing response", err)
		return
	}

	return repo.gateway.CreateResource(repo.config.ApiEndpoint(), path, bytes.NewReader(jsonBytes))
}

func (repo CCUserProvidedServiceInstanceRepository) Update(serviceInstanceFields models.ServiceInstanceFields) (apiErr error) {
	path := fmt.Sprintf("/v2/user_provided_service_instances/%s", serviceInstanceFields.Guid)

	type RequestBody struct {
		Credentials    map[string]interface{} `json:"credentials,omitempty"`
		SysLogDrainUrl string                 `json:"syslog_drain_url,omitempty"`
	}

	reqBody := RequestBody{serviceInstanceFields.Params, serviceInstanceFields.SysLogDrainUrl}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		apiErr = errors.NewWithError("Error parsing response", err)
		return
	}

	return repo.gateway.UpdateResource(repo.config.ApiEndpoint(), path, bytes.NewReader(jsonBytes))
}
