package docker

import "errors"

func Health() (bool, error) {
	// make sure docker can be run
	return false, errors.New("TODO")
}

func StartContainer() error {
	return nil
}

func BuildImage() error {
	return nil
}

func DeleteContainer() error {
	return nil
}

func StopContainer() error {
	return nil
}
