package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
)

// imagesReference allows build systems to inject imagesReference for CCCMO components
type imagesReference struct {
	CloudControllerManagerAWS       string `json:"cloudControllerManagerAWS"`
	CloudControllerManagerOpenStack string `json:"cloudControllerManagerOpenStack"`
}

// operatorConfig contains configuration values for templating resources
type operatorConfig struct {
	ManagedNamespace string
	ControllerImage  string
}

func getProviderFromInfrastructure(infra *configv1.Infrastructure) (configv1.PlatformType, error) {
	if infra == nil || infra.Status.PlatformStatus == nil {
		return "", fmt.Errorf("platform status is not pupulated on infrastructure")
	}
	if infra.Status.PlatformStatus.Type == "" {
		return "", fmt.Errorf("no platform provider found on infrastructure")
	}

	return infra.Status.PlatformStatus.Type, nil
}

func getImagesFromJSONFile(filePath string) (imagesReference, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return imagesReference{}, err
	}

	i := imagesReference{}
	if err := json.Unmarshal(data, &i); err != nil {
		return imagesReference{}, err
	}
	return i, nil
}

func getProviderControllerFromImages(platform configv1.PlatformType, images imagesReference) string {
	switch platform {
	case configv1.AWSPlatformType:
		return images.CloudControllerManagerAWS
	case configv1.OpenStackPlatformType:
		return images.CloudControllerManagerOpenStack
	default:
		return ""
	}
}