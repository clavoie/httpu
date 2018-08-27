package httpu_test

import "errors"

type ErrType int

func (et ErrType) MarshalJSON() ([]byte, error) {
	return nil, errors.New("error")
}

func (et *ErrType) UnmarshalJSON(data []byte) error {
	return errors.New("error")
}
