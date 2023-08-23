package main

import (
	"fmt"

	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type AddAppOptions struct {
	App string `positional-arg-name:"ORG" required:"true" description:"App to resume"`
	Key string `positional-arg-name:"KEY" required:"true" description:"Key of the annotation or label"`
	Val string `positional-arg-name:"VALUE" required:"true" description:"Value of the annotation or label"`
}

type AddAppAnnotationCommand struct {
	API           string        `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	AddAppOptions AddAppOptions `required:"2" positional-args:"true"`
}

type AddAppLabelCommand struct {
	API           string        `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	AddAppOptions AddAppOptions `required:"2" positional-args:"true"`
}

var addAppAnnotationCommand AddAppAnnotationCommand
var addAppLabelCommand AddAppLabelCommand

func ExecuteAddingApp(elementType MetadataElementType, appName string, key string, value string) error {
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

	space, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return err
	}
	spaceShow := space.Name
	if spaceShow != "" {
		spaceShow = fmt.Sprint(messages.C.Cyan(spaceShow))
	}

	appShow := appName
	if appShow != "" {
		appShow = fmt.Sprint(messages.C.Cyan(appShow))
	}
	_, _ = messages.Printf("Adding %s %s to app %s in space %s in org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), appShow, spaceShow, orgShow, messages.C.Cyan(username))

	appID, err := getAppID(appName)
	if err != nil {
		return err
	}

	data := &CommonResource{}
	data.AddMetadataElement(elementType, key, value)
	err = patchResource(AppType, appID, data)
	if err != nil {
		messages.Errorf("Error adding %s %s to app %s: %s\n", elementType, key, appName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func (c *AddAppAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingApp(MetadataAnnotationType, c.AddAppOptions.App, c.AddAppOptions.Key, c.AddAppOptions.Val)
}

func (c *AddAppLabelCommand) Execute(_ []string) error {
	return ExecuteAddingApp(MetadataLabelType, c.AddAppOptions.App, c.AddAppOptions.Key, c.AddAppOptions.Val)
}

func init() {
	desc := `Add an annotation to an app.`
	_, err := parser.AddCommand(
		"add-app-annotation",
		desc,
		desc,
		&addAppAnnotationCommand)
	if err != nil {
		panic(err)
	}
	desc = `Add a label to an app.`
	_, err = parser.AddCommand(
		"add-app-label",
		desc,
		desc,
		&addAppLabelCommand)
	if err != nil {
		panic(err)
	}
}
