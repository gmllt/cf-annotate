package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/types"
	"github.com/gmllt/cf-annotate/utils"
	"github.com/gmllt/cf-annotate/utils/messages"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gmllt/cf-annotate/metadata"
)

type AddSpaceOptions struct {
	Space string `positional-arg-name:"SPACE" required:"true" description:"Space to set metadata to"`
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

type ListSpaceOptions struct {
	Space string `positional-arg-name:"SPACE" required:"true" description:"Space to read"`
}

type ListSpaceAnnotationCommand struct {
	API              string           `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListSpaceOptions ListSpaceOptions `required:"2" positional-args:"true"`
}

type ListSpaceLabelCommand struct {
	API              string           `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListSpaceOptions ListSpaceOptions `required:"2" positional-args:"true"`
}

type RemoveSpaceOptions struct {
	Space string `positional-arg-name:"SPACE" required:"true" description:"Space to remove metadata from"`
	Key   string `positional-arg-name:"KEY" required:"true" description:"Key of the annotation or label"`
}

type RemoveSpaceAnnotationCommand struct {
	API                string             `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	RemoveSpaceOptions RemoveSpaceOptions `required:"2" positional-args:"true"`
}

type RemoveSpaceLabelCommand struct {
	API                string             `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	RemoveSpaceOptions RemoveSpaceOptions `required:"2" positional-args:"true"`
}

func ExecuteAddingSpace(elementType metadata.ElementType, spaceName string, key string, value string) error {
	username, err := utils.CliConnection.Username()
	if err != nil {
		return err
	}

	org, err := utils.CliConnection.GetCurrentOrg()
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

	spaceID, err := utils.GetID(utils.SpaceType, spaceName)
	if err != nil {
		return err
	}

	data := &metadata.CommonResource{}
	data.AddMetadataElement(elementType, key, value)
	err = utils.PatchResource(utils.SpaceType, spaceID, data)
	if err != nil {
		messages.Errorf("Error adding %s %s to space %s: %s\n", elementType, key, spaceName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func ExecuteListSpace(elementType metadata.ElementType, spaceName string) error {
	username, err := utils.CliConnection.Username()
	if err != nil {
		return err
	}

	org, err := utils.CliConnection.GetCurrentOrg()
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
	_, _ = messages.Printf("Listing %s of space %s in org %s as %s\n", messages.C.Cyan(fmt.Sprintf("%ss", elementType)), spaceShow, orgShow, messages.C.Cyan(username))

	spaceID, err := utils.GetID(utils.SpaceType, spaceName)
	if err != nil {
		return err
	}

	data, err := utils.GetResource(utils.SpaceType, spaceID)
	if err != nil {
		return err
	}

	elements := make(map[string]types.NullString)

	if elementType == metadata.AnnotationType {
		for k, v := range data.Metadata.Annotations {
			elements[k] = types.NullString{
				Value: v.Value,
				IsSet: v.IsSet,
			}
		}
	} else {
		for k, v := range data.Metadata.Labels {
			elements[k] = types.NullString{
				Value: v.Value,
				IsSet: v.IsSet,
			}
		}
	}

	messages.PrintMetadata(elements)
	return nil
}

func ExecuteRemovingSpace(elementType metadata.ElementType, spaceName string, key string) error {
	username, err := utils.CliConnection.Username()
	if err != nil {
		return err
	}

	org, err := utils.CliConnection.GetCurrentOrg()
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
	_, _ = messages.Printf("Removing %s %s from space %s in org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), spaceShow, orgShow, messages.C.Cyan(username))

	spaceID, err := utils.GetID(utils.SpaceType, spaceName)
	if err != nil {
		return err
	}

	data := &metadata.CommonResource{}
	data.RemoveMetadataElement(elementType, key)
	err = utils.PatchResource(utils.AppType, spaceID, data)
	if err != nil {
		messages.Errorf("Error removing %s %s from space %s: %s\n", elementType, key, spaceName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func (c *RemoveSpaceAnnotationCommand) Execute(_ []string) error {
	return ExecuteRemovingSpace(metadata.AnnotationType, c.RemoveSpaceOptions.Space, c.RemoveSpaceOptions.Key)
}

func (c *RemoveSpaceLabelCommand) Execute(_ []string) error {
	return ExecuteRemovingSpace(metadata.LabelType, c.RemoveSpaceOptions.Space, c.RemoveSpaceOptions.Key)
}

func (c *ListSpaceAnnotationCommand) Execute(_ []string) error {
	return ExecuteListSpace(metadata.AnnotationType, c.ListSpaceOptions.Space)
}

func (c *ListSpaceLabelCommand) Execute(_ []string) error {
	return ExecuteListSpace(metadata.LabelType, c.ListSpaceOptions.Space)
}

func (c *AddSpaceAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingSpace(metadata.AnnotationType, c.AddSpaceOptions.Space, c.AddSpaceOptions.Key, c.AddSpaceOptions.Val)
}

func (c *AddSpaceLabelCommand) Execute(_ []string) error {
	return ExecuteAddingSpace(metadata.LabelType, c.AddSpaceOptions.Space, c.AddSpaceOptions.Key, c.AddSpaceOptions.Val)
}

func init() {
	desc := fmt.Sprintf("%s a %s to a %s.", cases.Title(language.English, cases.Compact).String(addCommand), annotationElement, spaceResource)
	_, err := parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", addCommand, spaceResource, annotationElement),
		desc,
		desc,
		&AddOrgAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s to a %s.", cases.Title(language.English, cases.Compact).String(addCommand), labelElement, spaceResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", addCommand, spaceResource, labelElement),
		desc,
		desc,
		&AddSpaceLabelCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s all %ss of a %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(listCommand), annotationElement, spaceResource, spaceResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", listCommand, spaceResource, annotationElement),
		desc,
		desc,
		&ListSpaceAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s all %ss of a %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(listCommand), labelElement, spaceResource, spaceResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", listCommand, spaceResource, labelElement),
		desc,
		desc,
		&ListSpaceLabelCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s from a %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), annotationElement, spaceResource, annotationElement)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", removeCommand, spaceResource, annotationElement),
		desc,
		desc,
		&RemoveSpaceAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s from a %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), labelElement, spaceResource, labelElement)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", removeCommand, spaceResource, labelElement),
		desc,
		desc,
		&RemoveSpaceLabelCommand{})
	if err != nil {
		panic(err)
	}
}
