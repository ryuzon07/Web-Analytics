package main

import (
	"context"
	"embed" // <-- Keep
	"io/fs"
	"log"
	"net/http" // <-- Keep

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"

	"go-analytics/pkg/config"
	db "go-analytics/db/sqlc"
)
//go:embed ui/static
var staticFiles embed.FS

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBrokerURLs...),
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	handler := NewHandler(kafkaWriter, queries)

	// Initialize Gin router
	r := gin.Default()

	// --- UI Routes ---
	// Create a sub filesystem rooted at "ui/static" so the handler can open files like "index.html"
	staticFS, err := fs.Sub(staticFiles, "ui/static")
	if err != nil {
		log.Fatalf("could not create static sub-filesystem: %v", err)
	}

	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/", func(c *gin.Context) {
		// Open 'index.html' from the sub filesystem
		file, err := staticFS.Open("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "could not open index.html")
			return
		}
		defer file.Close()
		
		stat, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "could not stat index.html")
			return
		}

		c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", file, nil)
	})
	// --- END: UI Routes ---

	// API Routes
	r.POST("/events", handler.PostEvent)
	r.GET("/stats", handler.GetStats)

	log.Println("API service (with Web UI) starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}