# cf-annotate

A simple plugin to help you to annotate or label your Cloud Foundry organizations, spaces and applications.

## Installation

1. Download latest release made for your os from [releases page](https://github.com/gmllt/cf-annotate/releases)
2. Run `cf install-plugin <path-to-downloaded-binary>`

## Usage

```
    set-org-label ORG_NAME KEY VALUE        Add a label to an org
    set-space-label SPACE_NAME KEY VALUE    Add a label to a space
    set-app-label APP_NAME KEY VALUE        Add a label to an app
```

## Planned

```
    delete-org-label ORG_NAME KEY           Remove a label from an org
    delete-space-label SPACE_NAME KEY       Remove a label from a space
    delete-app-label APP_NAME KEY           Remove a label from an app
    show-org-labels ORG_NAME                List labels of an org
    show-space-labels SPACE_NAME            List labels of a space
    show-app-labels APP_NAME                List labels of an app
```
