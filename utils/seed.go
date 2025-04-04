package utils

import (
	"context"
	"database/sql"
	"delivery-service/models"
	"fmt"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func insertCampaign(db *sql.DB, campaign models.Campaign) error {
	query := `INSERT INTO campaigns (campaign_id, image_creative, cta, state) VALUES ($1, $2, $3, $4) ON CONFLICT (campaign_id) DO NOTHING`
	_, err := db.Exec(query, campaign.CampaignID, campaign.ImageCreative, campaign.CTA, campaign.State)
	if err != nil {
		return err
	}
	query = `INSERT INTO rules (campaign_id, include_country, exclude_country, include_os, exclude_os, include_apps) 
	         VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, campaign.CampaignID,
		pq.Array(campaign.Rules.IncludeCountry), pq.Array(campaign.Rules.ExcludeCountry),
		pq.Array(campaign.Rules.IncludeOS), pq.Array(campaign.Rules.ExcludeOS),
		pq.Array(campaign.Rules.IncludeApps),
	)
	if err != nil {
		return err
	}

	fmt.Println("Campaign and Rules inserted successfully!")
	return nil
}
func insertDB() {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	campaigns := []models.Campaign{
		{
			CampaignID:    "spotify",
			ImageCreative: "https://somelink",
			CTA:           "Download",
			State:         "ACTIVE",
			Rules: &models.TargetingRules{
				IncludeCountry: []string{"US", "Canada"},
			},
		},
		{
			CampaignID:    "duolingo",
			ImageCreative: "https://somelink2",
			CTA:           "Install",
			State:         "ACTIVE",
			Rules: &models.TargetingRules{
				IncludeOS:      []string{"Android", "iOS"},
				ExcludeCountry: []string{"US"},
			},
		},
		{
			CampaignID:    "subwaysurfer",
			ImageCreative: "https://somelink3",
			CTA:           "Play",
			State:         "ACTIVE",
			Rules: &models.TargetingRules{
				IncludeOS:   []string{"Android"},
				IncludeApps: []string{"com.gametion.ludokinggame"},
			},
		},
	}

	for _, campaign := range campaigns {
		err := insertCampaign(db, campaign)
		if err != nil {
			log.Println("error", err)
		}
	}
}
func UpdateCampaignState(db *sql.DB, campaignID string, newState string) error {
	ctx := context.Background()

	_, err := db.Exec("UPDATE campaigns SET state = $1 WHERE campaign_id = $2", newState, campaignID)
	if err != nil {
		return err
	}
	campaigns, err := GetCampaignsFromDB()
	if err != nil {
		return err
	}

	err = updateRedisCache(ctx, campaigns)
	if err != nil {
		fmt.Println("Failed to update Redis cache:", err)
	}

	return nil
}
