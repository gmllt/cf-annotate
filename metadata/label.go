package metadata

import "code.cloudfoundry.org/cli/types"

type Label types.NullString

func (l *Label) UnmarshalJSON(data []byte) error {
	var s types.NullString
	err := s.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	*l = Label(s)
	return nil
}

func (l Label) MarshalJSON() ([]byte, error) {
	return types.NullString(l).MarshalJSON()
}
