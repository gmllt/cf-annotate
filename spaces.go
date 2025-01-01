package main

import (
	"fmt"
	"github.com/gmllt/cf-annotate/metadata"

	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type AddSpaceOptions struct {
	Space string `positional-arg-name:"SPACE" required:"true" description:"Space to resume"`
	Key   string `positional-arg-name:"KEY" required:"true" description:"Key of the annotation or label"`
	Val   string `positional-arg-name:"VALUE" required:"true" description:"Value of the annotation or label"`
}

type AddSpaceAnnotationCommand struct {
	API             string          `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	AddSpaceOptions AddSpaceOptions `required:"2" positional-args:"true"`
}

type AddSpaceLabelCommand struct {
	API             string          `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	AddSpaceOptions AddSpaceOptions `required:"2" positional-args:"true"`
}

var addSpaceAnnotationCommand AddSpaceAnnotationCommand
var addSpaceLabelCommand AddSpaceLabelCommand

func ExecuteAddingSpace(elementType metadata.MetadataElementType, spaceName string, key string, value string) error {
	username, err := cliConnection.Username()
	if err != nil {
		return err
	}

	org, err := cliConnection.GetCurrentOrg()
	if err != nil {
		return err
	}
	orgShow := org.Name
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}

	spaceShow := spaceName
	if spaceShow != "" {
		spaceShow = fmt.Sprint(messages.C.Cyan(spaceShow))
	}
	_, _ = messages.Printf("Adding %s %s to space %s in org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), spaceShow, orgShow, messages.C.Cyan(username))

	spaceID, err := getSpaceID(spaceName)
	if err != nil {
		return err
	}

	data := &CommonResource{}
	data.AddMetadataElement(elementType, key, value)
	err = patchResource(SpaceType, spaceID, data)
	if err != nil {
		messages.Errorf("Error adding %s %s to space %s: %s\n", elementType, key, spaceName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func (c *AddSpaceAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingSpace(metadata.MetadataAnnotationType, c.AddSpaceOptions.Space, c.AddSpaceOptions.Key, c.AddSpaceOptions.Val)
}

func (c *AddSpaceLabelCommand) Execute(_ []string) error {
	return ExecuteAddingSpace(metadata.MetadataLabelType, c.AddSpaceOptions.Space, c.AddSpaceOptions.Key, c.AddSpaceOptions.Val)
}

func init() {
	desc := `Add an annotation to an space.`
	_, err := parser.AddCommand(
		"add-space-annotation",
		desc,
		desc,
		&addSpaceAnnotationCommand)
	if err != nil {
		panic(err)
	}
	desc = `Add a label to an space.`
	_, err = parser.AddCommand(
		"add-space-label",
		desc,
		desc,
		&addSpaceLabelCommand)
	if err != nil {
		panic(err)
	}
}
