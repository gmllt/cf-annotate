package main

import (
	"fmt"

	"github.com/gmllt/cf-annotate/metadata"

	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type AddOrgOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to resume"`
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
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to resume"`
}

type ListOrgAnnotationCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListOrgOptions ListOrgOptions `required:"2" positional-args:"true"`
}

type ListOrgLabelCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListOrgOptions ListOrgOptions `required:"2" positional-args:"true"`
}

var addOrgAnnotationCommand AddOrgAnnotationCommand
var addOrgLabelCommand AddOrgLabelCommand

func ExecuteAddingOrg(elementType metadata.MetadataElementType, orgName string, key string, value string) error {
	username, err := cliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := orgName
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Adding %s %s to org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), orgShow, messages.C.Cyan(username))

	orgID, err := getOrgID(orgName)
	if err != nil {
		return err
	}

	data := &CommonResource{}
	data.AddMetadataElement(elementType, key, value)
	err = patchResource(OrgType, orgID, data)
	if err != nil {
		messages.Errorf("Error adding %s %s to org %s: %s\n", elementType, key, orgName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func ExecuteListOrg(elementType metadata.MetadataElementType, orgName string) error {
	orgID, err := getOrgID(orgName)
	if err != nil {
		return err
	}

	data, err := getOrg(orgID)
	if err != nil {
		return err
	}

	if elementType == metadata.MetadataAnnotationType {
		for k, v := range data.Metadata.Annotations {
			_, _ = messages.Printf("%s: %s\n", k, v)
		}
	} else {
		for k, v := range data.Metadata.Labels {
			_, _ = messages.Printf("%s: %s\n", k, v)
		}
	}
	return nil
}

func (c *ListOrgAnnotationCommand) Execute(_ []string) error {
	return ExecuteListOrg(metadata.MetadataAnnotationType, c.ListOrgOptions.Org)
}

func (c *ListOrgLabelCommand) Execute(_ []string) error {
	return ExecuteListOrg(metadata.MetadataLabelType, c.ListOrgOptions.Org)
}

func (c *AddOrgAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingOrg(metadata.MetadataAnnotationType, c.AddOrgOptions.Org, c.AddOrgOptions.Key, c.AddOrgOptions.Val)
}

func (c *AddOrgLabelCommand) Execute(_ []string) error {
	return ExecuteAddingOrg(metadata.MetadataLabelType, c.AddOrgOptions.Org, c.AddOrgOptions.Key, c.AddOrgOptions.Val)
}

func init() {
	desc := `Add an annotation to an organization.`
	_, err := parser.AddCommand(
		"add-org-annotation",
		desc,
		desc,
		&addOrgAnnotationCommand)
	if err != nil {
		panic(err)
	}
	desc = `Add a label to an organization.`
	_, err = parser.AddCommand(
		"add-org-label",
		desc,
		desc,
		&addOrgLabelCommand)
	if err != nil {
		panic(err)
	}
	desc = `List all annotations of an organization.`
	_, err = parser.AddCommand(
		"list-org-annotation",
		desc,
		desc,
		&ListOrgAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = `List all labels of an organization.`
	_, err = parser.AddCommand(
		"list-org-label",
		desc,
		desc,
		&ListOrgLabelCommand{})
	if err != nil {
		panic(err)
	}
}
