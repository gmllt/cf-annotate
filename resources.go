package main

import (
	"github.com/gmllt/cf-annotate/metadata"
)

type CommonResource struct {
	Metadata metadata.Metadata `json:"metadata"`
}

type GUIDResource struct {
	CommonResource
	GUID string `json:"guid"`
}

func (r *CommonResource) AddMetadataElement(element metadata.MetadataElementType, key string, value string) {
	switch element {
	case metadata.MetadataAnnotationType:
		if r.Metadata.Annotations == nil {
			r.Metadata.Annotations = make(map[string]metadata.Annotation)
		}
		r.Metadata.Annotations[key] = metadata.Annotation{Value: value, IsSet: true}
	case metadata.MetadataLabelType:
		if r.Metadata.Labels == nil {
			r.Metadata.Labels = make(map[string]metadata.Label)
		}
		r.Metadata.Labels[key] = metadata.Label{Value: value, IsSet: true}
	}
}

func (r *CommonResource) RemoveMetadataElement(element metadata.MetadataElementType, key string) {
	r.AddMetadataElement(element, key, "")
}
