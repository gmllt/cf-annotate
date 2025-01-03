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

type AddOrgOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to set metadata to"`
	Key string `positional-arg-name:"KEY" required:"true" description:"Key of the annotation or label"`
	Val string `positional-arg-name:"VALUE" required:"true" description:"Value of the annotation or label"`
}

type AddOrgAnnotationCommand struct {
	API           string        `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	AddOrgOptions AddOrgOptions `required:"2" positional-args:"true"`
}

type AddOrgLabelCommand struct {
	API           string        `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	AddOrgOptions AddOrgOptions `required:"2" positional-args:"true"`
}

type ListOrgOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to read"`
}

type ListOrgAnnotationCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListOrgOptions ListOrgOptions `required:"2" positional-args:"true"`
}

type ListOrgLabelCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListOrgOptions ListOrgOptions `required:"2" positional-args:"true"`
}

type RemoveOrgOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to remove metadata from"`
	Key string `positional-arg-name:"KEY" required:"true" description:"Key of the annotation or label"`
}

type RemoveOrgAnnotationCommand struct {
	API              string           `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	RemoveOrgOptions RemoveOrgOptions `required:"2" positional-args:"true"`
}

type RemoveOrgLabelCommand struct {
	API              string           `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	RemoveOrgOptions RemoveOrgOptions `required:"2" positional-args:"true"`
}

var addOrgAnnotationCommand AddOrgAnnotationCommand
var addOrgLabelCommand AddOrgLabelCommand

func ExecuteAddingOrg(elementType metadata.ElementType, orgName string, key string, value string) error {
	username, err := utils.CliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := orgName
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Adding %s %s to org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), orgShow, messages.C.Cyan(username))

	orgID, err := utils.GetID(utils.OrgType, orgName)
	if err != nil {
		return err
	}

	data := &metadata.CommonResource{}
	data.AddMetadataElement(elementType, key, value)
	err = utils.PatchResource(utils.OrgType, orgID, data)
	if err != nil {
		messages.Errorf("Error adding %s %s to org %s: %s\n", elementType, key, orgName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func ExecuteListOrg(elementType metadata.ElementType, orgName string) error {
	username, err := utils.CliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := orgName
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Listing %s of org %s as %s\n", messages.C.Cyan(fmt.Sprintf("%ss", elementType)), orgShow, messages.C.Cyan(username))

	orgID, err := utils.GetID(utils.OrgType, orgName)
	if err != nil {
		return err
	}

	data, err := utils.GetResource(utils.OrgType, orgID)
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

func ExecuteRemovingOrg(elementType metadata.ElementType, orgName string, key string) error {
	username, err := utils.CliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := orgName
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Removing %s %s from org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), orgShow, messages.C.Cyan(username))

	orgID, err := utils.GetID(utils.OrgType, orgName)
	if err != nil {
		return err
	}

	data := &metadata.CommonResource{}
	data.RemoveMetadataElement(elementType, key)
	err = utils.PatchResource(utils.OrgType, orgID, data)
	if err != nil {
		messages.Errorf("Error removing %s %s from org %s: %s\n", elementType, key, orgName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func (c *RemoveOrgAnnotationCommand) Execute(_ []string) error {
	return ExecuteRemovingOrg(metadata.AnnotationType, c.RemoveOrgOptions.Org, c.RemoveOrgOptions.Key)
}

func (c *RemoveOrgLabelCommand) Execute(_ []string) error {
	return ExecuteRemovingOrg(metadata.LabelType, c.RemoveOrgOptions.Org, c.RemoveOrgOptions.Key)
}

func (c *ListOrgAnnotationCommand) Execute(_ []string) error {
	return ExecuteListOrg(metadata.AnnotationType, c.ListOrgOptions.Org)
}

func (c *ListOrgLabelCommand) Execute(_ []string) error {
	return ExecuteListOrg(metadata.LabelType, c.ListOrgOptions.Org)
}

func (c *AddOrgAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingOrg(metadata.AnnotationType, c.AddOrgOptions.Org, c.AddOrgOptions.Key, c.AddOrgOptions.Val)
}

func (c *AddOrgLabelCommand) Execute(_ []string) error {
	return ExecuteAddingOrg(metadata.LabelType, c.AddOrgOptions.Org, c.AddOrgOptions.Key, c.AddOrgOptions.Val)
}

func init() {
	desc := fmt.Sprintf("%s a %s to an %s.", cases.Title(language.English, cases.Compact).String(addCommand), annotationElement, organizationResource)
	_, err := parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", addCommand, organizationResource, annotationElement),
		desc,
		desc,
		&addOrgAnnotationCommand)
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s to an %s.", cases.Title(language.English, cases.Compact).String(addCommand), labelElement, organizationResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", addCommand, organizationResource, labelElement),
		desc,
		desc,
		&addOrgLabelCommand)
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s all %ss of an %s.", cases.Title(language.English, cases.Compact).String(listCommand), annotationElement, organizationResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s", listCommand, organizationResource),
		desc,
		desc,
		&ListOrgAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s all %ss of an %s.", cases.Title(language.English, cases.Compact).String(listCommand), labelElement, organizationResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s", listCommand, organizationResource),
		desc,
		desc,
		&ListOrgLabelCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s from an %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), annotationElement, organizationResource, annotationElement)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", removeCommand, organizationResource, annotationElement),
		desc,
		desc,
		&RemoveOrgAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s from an %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), labelElement, organizationResource, labelElement)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", removeCommand, organizationResource, labelElement),
		desc,
		desc,
		&RemoveOrgLabelCommand{})
	if err != nil {
		panic(err)
	}
}
