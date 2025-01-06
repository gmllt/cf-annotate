package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gmllt/cf-annotate/utils"
	"github.com/gmllt/cf-annotate/utils/messages"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/jessevdk/go-flags"
)

type Options struct {
}

var (
	pluginVersion = "0.1.3"
	options       Options
	parser        = flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)

	organizationResource = "org"
	spaceResource        = "space"
	appResource          = "app"
	annotationElement    = "annotation"
	labelElement         = "label"
	addCommand           = "set"
	removeCommand        = "delete"
	listCommand          = "show"
)

func Parse(args []string) error {
	_, err := parser.ParseArgs(args)
	if err != nil {
		var errFlag *flags.Error
		if errors.As(err, &errFlag) && errFlag.Type == flags.ErrCommandRequired {
			return nil
		}
		if errors.As(err, &errFlag) && errFlag.Type == flags.ErrHelp {
			messages.Errorf("Error parsing arguments: %s", err)
			return nil
		}
		return err
	}

	return nil
}

type AnnotatePlugin struct {
	Connection plugin.CliConnection
	Out        io.Writer
}

func (c *AnnotatePlugin) GetMetadata() plugin.PluginMetadata {
	var major, minor, build int
	_, _ = fmt.Sscanf(pluginVersion, "%d.%d.%d", &major, &minor, &build)

	// Generate commands
	var commands []plugin.Command
	for _, resource := range []string{organizationResource, spaceResource, appResource} {
		for _, element := range []string{annotationElement, labelElement} {
			commands = append(commands, plugin.Command{
				Name:     fmt.Sprintf("%s-%s-%s", addCommand, resource, element),
				HelpText: fmt.Sprintf("%s %s to a %s.", cases.Title(language.English, cases.Compact).String(addCommand), element, resource),
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s-%s-%s %s_NAME KEY VALUE", addCommand, resource, element, strings.ToUpper(resource)),
				},
			})
			commands = append(commands, plugin.Command{
				Name:     fmt.Sprintf("%s-%s-%s", removeCommand, resource, element),
				HelpText: fmt.Sprintf("%s %s from a %s. If the %s does not exist, nothing happens.", cases.Title(language.English, cases.Compact).String(removeCommand), element, resource, element),
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s-%s-%s %s_NAME KEY", removeCommand, resource, element, strings.ToUpper(resource)),
				},
			})

			commands = append(commands, plugin.Command{
				Name:     fmt.Sprintf("%s-%s-%s", listCommand, resource, element),
				HelpText: fmt.Sprintf("%s all %ss of a %s.", cases.Title(language.English, cases.Compact).String(listCommand), element, resource),
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s-%s-%s %s_NAME", listCommand, resource, element, strings.ToUpper(resource)),
				},
			})
		}
	}
	return plugin.PluginMetadata{
		Name: "AnnotatePlugin",
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: build,
		},
		Commands: commands,
	}
}

func (c *AnnotatePlugin) Run(cc plugin.CliConnection, args []string) {
	utils.CliConnection = cc

	action := args[0]
	if action == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	err := Parse(args)
	if err != nil {
		messages.Fatal(err.Error())
	}
}
