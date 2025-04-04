package models

type Campaign struct {
	CampaignID    string
	CampaignName  string `json:"campaign_name,omitempty"`
	ImageCreative string
	CTA           string
	State         string          `json:"state,omitempty"`
	Rules         *TargetingRules `json:"rules,omitempty"`
}
