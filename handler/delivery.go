package handler

import (
	"context"
	"delivery-service/models"
	"delivery-service/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "delivery_requests_total",
			Help: "Total number of delivery requests",
		},
		[]string{"method"},
	)
	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "delivery_request_duration_seconds",
			Help:    "Histogram of request durations",
			Buckets: prometheus.DefBuckets,
		},
	)
)
var (
	errorRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "delivery_errors_total",
			Help: "Total number of delivery errors",
		},
		[]string{"error_type"},
	)
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(errorRequests)
}

func Gethandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	totalRequests.WithLabelValues(r.Method).Inc()
	defer func() {
		duration := time.Since(start).Seconds()
		requestDuration.Observe(duration)
	}()
	fmt.Println("inside get handler")
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method request", http.StatusMethodNotAllowed)
		return
	}
	appID := r.URL.Query().Get("app")
	country := r.URL.Query().Get("country")
	os := r.URL.Query().Get("os")
	fmt.Println(appID, country, os)
	if appID == "" {
		errorRequests.WithLabelValues("missing_app").Inc()
		http.Error(w, `{"error": "missing app param"}`, http.StatusBadRequest)
		return
	}
	if country == "" {

		errorRequests.WithLabelValues("missing_country").Inc()
		http.Error(w, `{"error": "missing country param"}`, http.StatusBadRequest)
		return
	}
	if os == "" {
		errorRequests.WithLabelValues("missing_os").Inc()
		fmt.Println("os missed")
		http.Error(w, `{"error": "missing os param"}`, http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	start = time.Now()

	cachedData, err := utils.RedisClient.Get(ctx, "active_campaigns").Result()
	elapsed := time.Since(start)
	fmt.Printf("Redis GET operation took: %v\n", elapsed)
	var campaigns []models.Campaign
	if err == nil {
		// Redis cache hit, unmarshal JSON data
		err := json.Unmarshal([]byte(cachedData), &campaigns)
		if err != nil {
			fmt.Println("Failed to unmarshal Redis data:", err)
		} else {
			fmt.Println("Fetched campaigns from Redis")
		}
	}
	if err != nil || len(campaigns) == 0 {
		fmt.Println("Fetching campaigns from DB")
		start := time.Now()

		campaigns, err = utils.GetCampaignsFromDB()

		elapsed := time.Since(start)
		fmt.Printf("Database fetch operation took: %v\n", elapsed)
		if err != nil {
			http.Error(w, `{"error": "failed to fetch campaigns"}`, http.StatusInternalServerError)
			return
		}

		// Store campaigns in Redis with a TTL of 10 minutes
		campaignsJSON, _ := json.Marshal(campaigns)
		utils.RedisClient.Set(ctx, "active_campaigns", campaignsJSON, 10*time.Minute)
	}

	matchingCampaigns := []models.Campaign{}
	//campaigns, err := getCampaignsFromDB()
	if err != nil {
		http.Error(w, `{"error": "failed to fetch campaigns"}`, http.StatusInternalServerError)
		return
	}
	for _, campaign := range campaigns {
		fmt.Println(campaign)
		if campaign.State != "ACTIVE" {
			continue
		}

		if !matchesRules(campaign.Rules, appID, country, os) {
			continue
		}

		matchingCampaigns = append(matchingCampaigns, models.Campaign{
			CampaignID:    campaign.CampaignID,
			ImageCreative: campaign.ImageCreative,
			CTA:           campaign.CTA,
		})

	}
	if len(matchingCampaigns) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(matchingCampaigns)

}
func matchesRules(rules *models.TargetingRules, appID, country, os string) bool {
	// Check country inclusion/exclusion rules
	if len(rules.IncludeCountry) > 0 && !contains(rules.IncludeCountry, country) {
		return false
	}
	if len(rules.ExcludeCountry) > 0 && contains(rules.ExcludeCountry, country) {
		return false
	}

	if len(rules.IncludeOS) > 0 && !contains(rules.IncludeOS, os) {
		return false
	}
	if len(rules.ExcludeOS) > 0 && contains(rules.ExcludeOS, os) {
		return false
	}

	if len(rules.IncludeApps) > 0 && !contains(rules.IncludeApps, appID) {
		return false
	}

	return true
}

func contains(list []string, item string) bool {
	for _, val := range list {
		if strings.EqualFold(val, item) { // Case insensitive match
			return true
		}
	}
	return false
}
