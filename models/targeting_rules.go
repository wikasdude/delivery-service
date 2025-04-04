package models

type TargetingRules struct {
	IncludeApps    []string `json:"includeapps,omitempty"`
	IncludeCountry []string `json:"includecountry,omitempty"`
	ExcludeCountry []string `json:"excludecountry,omitempty"`
	IncludeOS      []string `json:"includeos,omitempty"`
	ExcludeOS      []string `json:"excludeos,omitempty"`
}
