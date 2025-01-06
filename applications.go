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

type AddAppOptions struct {
	App string `positional-arg-name:"APP" required:"true" description:"App to resume"`
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

type ListAppOptions struct {
	App string `positional-arg-name:"APP" required:"true" description:"App to read"`
}

type ListAppMetadataCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListAppOptions ListAppOptions `required:"2" positional-args:"true"`
}

type ListAppAnnotationCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListAppOptions ListAppOptions `required:"2" positional-args:"true"`
}

type ListAppLabelCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ListAppOptions ListAppOptions `required:"2" positional-args:"true"`
}

type RemoveAppOptions struct {
	App string `positional-arg-name:"APP" required:"true" description:"App to remove metadata from"`
	Key string `positional-arg-name:"KEY" required:"true" description:"Key of the annotation or label"`
}

type RemoveAppAnnotationCommand struct {
	API              string           `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	RemoveAppOptions RemoveAppOptions `required:"2" positional-args:"true"`
}

type RemoveAppLabelCommand struct {
	API              string           `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	RemoveAppOptions RemoveAppOptions `required:"2" positional-args:"true"`
}

func ExecuteAddingApp(elementType metadata.ElementType, appName string, key string, value string) error {
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

	space, err := utils.CliConnection.GetCurrentSpace()
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

	appID, err := utils.GetID(utils.AppType, appName)
	if err != nil {
		return err
	}

	data := &metadata.CommonResource{}
	data.AddMetadataElement(elementType, key, value)
	err = utils.PatchResource(utils.AppType, appID, data)
	if err != nil {
		messages.Errorf("Error adding %s %s to app %s: %s\n", elementType, key, appName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func ExecuteListAllApp(appName string) error {
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

	space, err := utils.CliConnection.GetCurrentSpace()
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
	_, _ = messages.Printf("Listing %s of app %s in space %s, in org %s as %s\n", messages.C.Cyan(fmt.Sprintf("%ss", "metadata")), appShow, spaceShow, orgShow, messages.C.Cyan(username))

	appID, err := utils.GetID(utils.AppType, appName)
	if err != nil {
		return err
	}

	data, err := utils.GetResource(utils.AppType, appID)
	if err != nil {
		return err
	}

	elements := make(map[string]types.NullString)
	for k, v := range data.Metadata.Annotations {
		elements[k] = types.NullString{
			Value: v.Value,
			IsSet: v.IsSet,
		}
	}
	_, _ = messages.Printfln("%ss :", cases.Title(language.English, cases.Compact).String(annotationElement))
	messages.PrintMetadata(elements)
	elements = make(map[string]types.NullString)
	for k, v := range data.Metadata.Labels {
		elements[k] = types.NullString{
			Value: v.Value,
			IsSet: v.IsSet,
		}
	}
	_, _ = messages.Printfln("%ss :", cases.Title(language.English, cases.Compact).String(labelElement))
	messages.PrintMetadata(elements)

	return nil
}

func ExecuteListApp(elementType metadata.ElementType, appName string) error {
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

	space, err := utils.CliConnection.GetCurrentSpace()
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
	_, _ = messages.Printf("Listing %s of app %s in space %s, in org %s as %s\n", messages.C.Cyan(fmt.Sprintf("%ss", elementType)), appShow, spaceShow, orgShow, messages.C.Cyan(username))

	appID, err := utils.GetID(utils.AppType, appName)
	if err != nil {
		return err
	}

	data, err := utils.GetResource(utils.AppType, appID)
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

func ExecuteRemovingApp(elementType metadata.ElementType, appName string, key string) error {
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

	space, err := utils.CliConnection.GetCurrentSpace()
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
	_, _ = messages.Printf("Removing %s %s from app %s in space %s in org %s as %s\n", messages.C.Cyan(elementType), messages.C.Cyan(key), appShow, spaceShow, orgShow, messages.C.Cyan(username))

	appID, err := utils.GetID(utils.AppType, appName)
	if err != nil {
		return err
	}

	data := &metadata.CommonResource{}
	data.RemoveMetadataElement(elementType, key)
	err = utils.PatchResource(utils.AppType, appID, data)
	if err != nil {
		messages.Errorf("Error removing %s %s from org %s: %s\n", elementType, key, appName, err)
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK"))
	return nil
}

func (c *RemoveAppAnnotationCommand) Execute(_ []string) error {
	return ExecuteRemovingApp(metadata.AnnotationType, c.RemoveAppOptions.App, c.RemoveAppOptions.Key)
}

func (c *RemoveAppLabelCommand) Execute(_ []string) error {
	return ExecuteRemovingApp(metadata.LabelType, c.RemoveAppOptions.App, c.RemoveAppOptions.Key)
}

func (c *ListAppMetadataCommand) Execute(_ []string) error {
	return ExecuteListAllApp(c.ListAppOptions.App)
}

func (c *ListAppAnnotationCommand) Execute(_ []string) error {
	return ExecuteListApp(metadata.AnnotationType, c.ListAppOptions.App)
}

func (c *ListAppLabelCommand) Execute(_ []string) error {
	return ExecuteListApp(metadata.LabelType, c.ListAppOptions.App)
}

func (c *AddAppAnnotationCommand) Execute(_ []string) error {
	return ExecuteAddingApp(metadata.AnnotationType, c.AddAppOptions.App, c.AddAppOptions.Key, c.AddAppOptions.Val)
}

func (c *AddAppLabelCommand) Execute(_ []string) error {
	return ExecuteAddingApp(metadata.LabelType, c.AddAppOptions.App, c.AddAppOptions.Key, c.AddAppOptions.Val)
}

func init() {
	desc := fmt.Sprintf("%s an %s to an %s.", cases.Title(language.English, cases.Compact).String(addCommand), annotationElement, appResource)
	_, err := parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", addCommand, appResource, annotationElement),
		desc,
		desc,
		&AddAppAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s to an %s.", cases.Title(language.English, cases.Compact).String(addCommand), labelElement, appResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", addCommand, appResource, labelElement),
		desc,
		desc,
		&AddAppLabelCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("List all %ss of an %s.", cases.Title(language.English, cases.Compact).String(listCommand), "metadata", appResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", listCommand, appResource, "metadata"),
		desc,
		desc,
		&ListAppMetadataCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("List annotations of an %s.", appResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", listCommand, appResource, annotationElement),
		desc,
		desc,
		&ListAppAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("List labels of an %s.", appResource)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", listCommand, appResource, labelElement),
		desc,
		desc,
		&ListAppLabelCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s an %s from an %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), annotationElement, appResource, annotationElement)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", removeCommand, appResource, annotationElement),
		desc,
		desc,
		&RemoveAppAnnotationCommand{})
	if err != nil {
		panic(err)
	}
	desc = fmt.Sprintf("%s a %s from an %s.\n   If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), labelElement, appResource, labelElement)
	_, err = parser.AddCommand(
		fmt.Sprintf("%s-%s-%s", removeCommand, appResource, labelElement),
		desc,
		desc,
		&RemoveAppLabelCommand{})
	if err != nil {
		panic(err)
	}
}
