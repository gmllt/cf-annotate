package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ResourceType string

const (
	OrgType   ResourceType = "organizations"
	SpaceType ResourceType = "spaces"
	AppType   ResourceType = "apps"
)

func getOrgID(orgName string) (string, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/organizations?names="+orgName)
	if err != nil {
		return "", err
	}
	var resources struct {
		Resources []GUIDResource `json:"resources"`
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
	org, err := cliConnection.GetCurrentOrg()
	if err != nil {
		return "", err
	}

	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/spaces?names="+spaceName+"&organization_guids="+org.Guid)
	if err != nil {
		return "", err
	}
	var resources struct {
		Resources []GUIDResource `json:"resources"`
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
	space, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return "", err
	}

	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/apps?names="+appName+"&space_guids="+space.Guid)
	if err != nil {
		return "", err
	}
	var resources struct {
		Resources []GUIDResource `json:"resources"`
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

func patchResource(resourceType ResourceType, resourceGUID string, data *CommonResource) error {
	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println("curl", "-X", "PATCH", "/v3/"+string(resourceType)+"/"+resourceGUID, "-d", string(dataString))
	_, err = cliConnection.CliCommandWithoutTerminalOutput("curl", "-X", "PATCH", "/v3/"+string(resourceType)+"/"+resourceGUID, "-d", string(dataString))
	// print command
	if err != nil {
		return err
	}

	return nil
}

func getOrg(resourceGUID string) (*CommonResource, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/organizations/"+resourceGUID)
	if err != nil {
		return nil, err
	}
	var resource CommonResource
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func getSpace(resourceGUID string) (*CommonResource, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/spaces/"+resourceGUID)
	if err != nil {
		return nil, err
	}
	var resource CommonResource
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func getApp(resourceGUID string) (*CommonResource, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/apps/"+resourceGUID)
	if err != nil {
		return nil, err
	}
	var resource CommonResource
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
