package metadata

type CommonResource struct {
	Metadata Metadata `json:"metadata"`
}

type GUIDResource struct {
	CommonResource
	GUID string `json:"guid"`
}

func (r *CommonResource) AddMetadataElement(element ElementType, key string, value string) {
	switch element {
	case AnnotationType:
		if r.Metadata.Annotations == nil {
			r.Metadata.Annotations = make(map[string]Annotation)
		}
		r.Metadata.Annotations[key] = Annotation{Value: value, IsSet: true}
	case LabelType:
		if r.Metadata.Labels == nil {
			r.Metadata.Labels = make(map[string]Label)
		}
		r.Metadata.Labels[key] = Label{Value: value, IsSet: true}
	}
}

func (r *CommonResource) RemoveMetadataElement(element ElementType, key string) {
	switch element {
	case AnnotationType:
		if r.Metadata.Annotations == nil {
			r.Metadata.Annotations = make(map[string]Annotation)
		}
		r.Metadata.Annotations[key] = Annotation{Value: "", IsSet: false}
	case LabelType:
		if r.Metadata.Labels == nil {
			r.Metadata.Labels = make(map[string]Label)
		}
		r.Metadata.Labels[key] = Label{Value: "", IsSet: false}
	}
}
