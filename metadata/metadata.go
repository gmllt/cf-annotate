package metadata

type MetadataElementType string

const (
	MetadataAnnotationType MetadataElementType = "annotation"
	MetadataLabelType      MetadataElementType = "label"
)

type Metadata struct {
	Labels      map[string]Label      `json:"labels,omitempty"`
	Annotations map[string]Annotation `json:"annotations,omitempty"`
}
