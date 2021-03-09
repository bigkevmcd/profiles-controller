package httpapi

import (
	profilesv1alpha1 "github.com/bigkevmcd/profiles-controller/api/v1alpha1"
)

type searchResponse struct {
	Profiles []profilesv1alpha1.Profile `json:"profiles"`
}
