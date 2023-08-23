package main

import (
	"fmt"

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

var addOrgAnnotationCommand AddOrgAnnotationCommand
var addOrgLabelCommand AddOrgLabelCommand

func ExecuteAddingOrg(elementType MetadataElementType, orgName string, key string, value string) error {
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

func (c *AddOrgAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingOrg(MetadataAnnotationType, c.AddOrgOptions.Org, c.AddOrgOptions.Key, c.AddOrgOptions.Val)
}

func (c *AddOrgLabelCommand) Execute(_ []string) error {
	return ExecuteAddingOrg(MetadataLabelType, c.AddOrgOptions.Org, c.AddOrgOptions.Key, c.AddOrgOptions.Val)
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
}
