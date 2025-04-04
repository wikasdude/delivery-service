package utils

import (
	"context"
	"database/sql"
	"delivery-service/models"
	"fmt"

	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "campaign_db"
)

func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("✅ Connected to PostgreSQL")
	return db, nil
}

func GetCampaignsFromDB() ([]models.Campaign, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT c.campaign_id, c.image_creative, c.cta, c.state,
		COALESCE(r.include_country, '{}'),
		COALESCE(r.exclude_country, '{}'),
		COALESCE(r.include_os, '{}'),
		COALESCE(r.exclude_os, '{}'),
		COALESCE(r.include_apps, '{}')
		FROM campaigns c
		LEFT JOIN rules r ON c.campaign_id = r.campaign_id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := []models.Campaign{}
	for rows.Next() {
		var c models.Campaign
		var r models.TargetingRules
		err := rows.Scan(&c.CampaignID, &c.ImageCreative, &c.CTA, &c.State,
			pq.Array(&r.IncludeCountry), pq.Array(&r.ExcludeCountry),
			pq.Array(&r.IncludeOS), pq.Array(&r.ExcludeOS), pq.Array(&r.IncludeApps))
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		c.Rules = &r
		campaigns = append(campaigns, c)
	}
	err = updateRedisCache(context.Background(), campaigns)
	if err != nil {
		fmt.Println("Failed to update Redis cache:", err)
	}

	return campaigns, nil
}
