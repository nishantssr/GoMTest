package models

import "github.com/ottob/go-semver/semver"

type Case struct {
	ID           int64
	Message      string
	GuideVersion *semver.Version
}
