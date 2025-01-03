package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gmllt/cf-annotate/metadata"
)

type ResourceType string

const (
	OrgType   ResourceType = "organizations"
	SpaceType ResourceType = "spaces"
	AppType   ResourceType = "apps"
)

var (
	CliConnection plugin.CliConnection
)

func GetID(resourceType ResourceType, resourceName string) (string, error) {
	switch resourceType {
	case OrgType:
		return getOrgID(resourceName)
	case SpaceType:
		return getSpaceID(resourceName)
	case AppType:
		return getAppID(resourceName)
	default:
		return "", fmt.Errorf("unknown resource type %s", resourceType)
	}
}

func GetResource(resourceType ResourceType, resourceGUID string) (*metadata.CommonResource, error) {
	switch resourceType {
	case OrgType:
		return getOrg(resourceGUID)
	case SpaceType:
		return getSpace(resourceGUID)
	case AppType:
		return getApp(resourceGUID)
	default:
		return nil, fmt.Errorf("unknown resource type %s", resourceType)
	}
}

func getOrgID(orgName string) (string, error) {
	result, err := CliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/organizations?names="+orgName)
	if err != nil {
		return "", err
	}
	var resources struct {
		Resources []metadata.GUIDResource `json:"resources"`
	}
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resources)
	if err != nil {
		return "", err
	}
	if len(resources.Resources) == 0 {
		return "", fmt.Errorf("resource %s not found", orgName)
	}
	return resources.Resources[0].GUID, nil
}

func getSpaceID(spaceName string) (string, error) {
	org, err := CliConnection.GetCurrentOrg()
	if err != nil {
		return "", err
	}

	result, err := CliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/spaces?names="+spaceName+"&organization_guids="+org.Guid)
	if err != nil {
		return "", err
	}
	var resources struct {
		Resources []metadata.GUIDResource `json:"resources"`
	}
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resources)
	if err != nil {
		return "", err
	}
	if len(resources.Resources) == 0 {
		return "", fmt.Errorf("resource %s not found", spaceName)
	}
	return resources.Resources[0].GUID, nil
}

func getAppID(appName string) (string, error) {
	space, err := CliConnection.GetCurrentSpace()
	if err != nil {
		return "", err
	}

	result, err := CliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/apps?names="+appName+"&space_guids="+space.Guid)
	if err != nil {
		return "", err
	}
	var resources struct {
		Resources []metadata.GUIDResource `json:"resources"`
	}
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resources)
	if err != nil {
		return "", err
	}
	if len(resources.Resources) == 0 {
		return "", fmt.Errorf("resource %s not found", appName)
	}
	return resources.Resources[0].GUID, nil
}

func PatchResource(resourceType ResourceType, resourceGUID string, data *metadata.CommonResource) error {
	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = CliConnection.CliCommandWithoutTerminalOutput("curl", "-X", "PATCH", "/v3/"+string(resourceType)+"/"+resourceGUID, "-d", string(dataString))
	if err != nil {
		return err
	}

	return nil
}

func getOrg(resourceGUID string) (*metadata.CommonResource, error) {
	result, err := CliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/organizations/"+resourceGUID)
	if err != nil {
		return nil, err
	}
	var resource metadata.CommonResource
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func getSpace(resourceGUID string) (*metadata.CommonResource, error) {
	result, err := CliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/spaces/"+resourceGUID)
	if err != nil {
		return nil, err
	}
	var resource metadata.CommonResource
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func getApp(resourceGUID string) (*metadata.CommonResource, error) {
	result, err := CliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/apps/"+resourceGUID)
	if err != nil {
		return nil, err
	}
	var resource metadata.CommonResource
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
