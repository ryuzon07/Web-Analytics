package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/kafka-go"

	db "go-analytics/db/sqlc"
	"go-analytics/pkg/types"
)

type Handler struct {
	kafkaWriter *kafka.Writer
	dbStore     db.Querier
}

func NewHandler(kafkaWriter *kafka.Writer, dbStore db.Querier) *Handler {
	return &Handler{
		kafkaWriter: kafkaWriter,
		dbStore:     dbStore,
	}
}

// PostEvent handles a new event ingestion
func (h *Handler) PostEvent(c *gin.Context) {
	var event types.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize event"})
		return
	}

	// Write to Kafka
	err = h.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Value: eventBytes,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue event"})
		return
	}

	// Return 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{"message": "Event accepted"})
}

// GetStats handles retrieving aggregated stats
func (h *Handler) GetStats(c *gin.Context) {
	siteID := c.Query("site_id")
	dateStr := c.Query("date")

	if siteID == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id and date are required"})
		return
	}

	// Parse the date
	queryDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Calculate date range (the full day)
	startOfDay := queryDate.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 1. Get total views and unique users
	statsParams := db.GetSiteStatsParams{
		SiteID:    siteID,
		Timestamp: pgtype.Timestamptz{Time: startOfDay, Valid: true},
		Timestamp2: pgtype.Timestamptz{Time: endOfDay, Valid: true},
	}
	
	stats, err := h.dbStore.GetSiteStats(c.Request.Context(), statsParams)
	if err != nil {
		log.Printf("Error getting site stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stats"})
		return
	}

	// 2. Get top paths
	topPathsParams := db.GetTopPathsParams{
		SiteID:    siteID,
		Timestamp: pgtype.Timestamptz{Time: startOfDay, Valid: true},
		Timestamp2: pgtype.Timestamptz{Time: endOfDay, Valid: true},
	}
	
	topPaths, err := h.dbStore.GetTopPaths(c.Request.Context(), topPathsParams)
	if err != nil {
		log.Printf("Error getting top paths: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top paths"})
		return
	}

	// Build the response
	type TopPathResponse struct {
		Path  string `json:"path"`
		Views int64  `json:"views"`
	}
	
	responsePaths := make([]TopPathResponse, len(topPaths))
	for i, p := range topPaths {
		responsePaths[i] = TopPathResponse{
			Path:  p.Path,
			Views: p.Views,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"site_id":       siteID,
		"date":          dateStr,
		"total_views":   stats.TotalViews,
		"unique_users":  stats.UniqueUsers,
		"top_paths":     responsePaths,
	})
}