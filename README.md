# cf-annotate

A simple plugin to help you to annotate or label your Cloud Foundry organizations, spaces and applications.

## Installation

1. Download latest release made for your os from [releases page](https://github.com/gmllt/cf-annotate/releases)
2. Run `cf install-plugin <path-to-downloaded-binary>`

## Usage

```
    add-org-label ORG_NAME KEY VALUE        Add a label to an org
    add-space-label SPACE_NAME KEY VALUE    Add a label to a space
    add-app-label APP_NAME KEY VALUE        Add a label to an app
```

## Planned

```
    remove-org-label ORG_NAME KEY           Remove a label from an org
    remove-space-label SPACE_NAME KEY       Remove a label from a space
    remove-app-label APP_NAME KEY           Remove a label from an app
    list-org-labels ORG_NAME                List labels of an org
    list-space-labels SPACE_NAME            List labels of a space
    list-app-labels APP_NAME                List labels of an app
```
