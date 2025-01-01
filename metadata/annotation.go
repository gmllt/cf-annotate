package metadata

import "code.cloudfoundry.org/cli/types"

type Annotation types.NullString

func (a *Annotation) UnmarshalJSON(data []byte) error {
	var s types.NullString
	err := s.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	*a = Annotation(s)
	return nil
}

func (a Annotation) MarshalJSON() ([]byte, error) {
	return types.NullString(a).MarshalJSON()
}
