package metadata

type ElementType string

const (
	AnnotationType ElementType = "annotation"
	LabelType      ElementType = "label"
)

type Metadata struct {
	Labels      map[string]Label      `json:"labels,omitempty"`
	Annotations map[string]Annotation `json:"annotations,omitempty"`
}
