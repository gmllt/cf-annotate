package main

import "code.cloudfoundry.org/cli/types"

type MetadataElementType string

const (
	MetadataAnnotationType MetadataElementType = "annotation"
	MetadataLabelType      MetadataElementType = "label"
)

type Annotation types.NullString
type Label types.NullString

type Metadata struct {
	Labels      map[string]Label      `json:"labels"`
	Annotations map[string]Annotation `json:"annotations"`
}

type CommonResource struct {
	Metadata Metadata `json:"metadata"`
}

type GUIDResource struct {
	CommonResource
	GUID string `json:"guid"`
}

func (r *CommonResource) AddMetadataElement(element MetadataElementType, key string, value string) {
	switch element {
	case MetadataAnnotationType:
		if r.Metadata.Annotations == nil {
			r.Metadata.Annotations = make(map[string]Annotation)
		}
		r.Metadata.Annotations[key] = Annotation{Value: value, IsSet: true}
	case MetadataLabelType:
		if r.Metadata.Labels == nil {
			r.Metadata.Labels = make(map[string]Label)
		}
		r.Metadata.Labels[key] = Label{Value: value, IsSet: true}
	}
}

func (r *CommonResource) RemoveMetadataElement(element MetadataElementType, key string) {
	r.AddMetadataElement(element, key, "")
}
