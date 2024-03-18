package docker

import "errors"

func Health() (bool, error) {
	// make sure docker can be run
	return false, errors.New("TODO")
}
