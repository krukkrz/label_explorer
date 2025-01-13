package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/krukkrz/label_loader/discogs"
	"github.com/krukkrz/label_loader/mapping"
	"github.com/krukkrz/label_loader/releasescounts"
	"golang.org/x/time/rate"

	_ "github.com/lib/pq"
)

func main() {
	labelId := flag.Int("label-id", 2175451, "Label ID from DiscogsAPI")
	log.Println("Start fetching information for record label: ", *labelId)

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	resultChan, err := loadLabelReleases(labelId)
	if err != nil {
		log.Fatalf("error encountered while fetching releases details: %v", err)
	}

	releaseCounts, err := mapping.DetailsToReleaseCount(*labelId, resultChan)
	if err != nil {
		log.Fatalf("error encountered during mapping: %v", err)
	}

	storeReleaseCountsInDatabase(releaseCounts)
	log.Printf("label %d information successfully stored!", *labelId)
}

func storeReleaseCountsInDatabase(releaseCounts []releasescounts.ReleaseCount) {
	log.Println("storing release counts in database...")

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("user=%s password=%s dbname=labels sslmode=disable", dbUser, dbPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	releasescounts.InitTable(db)
	err = releasescounts.Save(db, releaseCounts)
	if err != nil {
		log.Fatalf("could not save label information: %v", err)
	}
}

func loadLabelReleases(labelId *int) (chan mapping.Details, error) {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	httpClient := &http.Client{}
	client := &discogs.Client{
		HttpClient: httpClient,
		Key:        apiKey,
		Secret:     apiSecret,
	}

	releases, err := discogs.Releases(client, fmt.Sprintf("https://api.discogs.com/labels/%d/releases", *labelId))
	if err != nil {
		return nil, fmt.Errorf("error ocurred while fetching releases from discogs: %v", err)
	}
	log.Println("releases fetched...")

	var wg sync.WaitGroup
	errChan := make(chan error, len(releases))
	resultChan := make(chan mapping.Details, len(releases))

	// Rate limit of Discogs API: https://www.discogs.com/developers#page:home,header:home-rate-limiting
	// We need to keep it below rate limit, so we don't hit 429 error code.
	const requestsPerSecond = 0.9

	limiter := rate.NewLimiter(requestsPerSecond, 1)
	ctx := context.Background()

	for _, release := range releases {
		err = limiter.Wait(ctx)
		if err != nil {
			return resultChan, fmt.Errorf("context cancelled or error occurred: %v", err)
		}
		go fetchReleaseDetails(client, release, &wg, resultChan, errChan)
	}

	wg.Wait()
	log.Println("done fetching each release details...")
	close(errChan)
	close(resultChan)
	for err = range errChan {
		log.Println("error occurred while fetching release details:", err)
	}
	return resultChan, err
}

func fetchReleaseDetails(
	client *discogs.Client,
	release discogs.Release,
	wg *sync.WaitGroup,
	resultChan chan mapping.Details,
	errChan chan error,
) {
	wg.Add(1)
	defer wg.Done()
	details, err := discogs.ReleaseDetails(client, release.ResourceUrl)
	if err != nil {
		errChan <- err
	}
	log.Println("done fetching release: ", release.Id)
	resultChan <- mapping.Details{
		Id:     details.Id,
		Artist: release.Artist,
		Genres: details.Genres,
		Styles: details.Genres,
	}
}
