// Package main contains the main HTTP server code
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gorpl/internal/api"
	"gorpl/internal/database"
)

const (
	// URL for downloading XML file from the medicinal products registry
	xmlURL = "https://rejestry.ezdrowie.gov.pl/api/rpl/medicinal-products/public-pl-report/6.0.0/overall.xml"
	// API version
	apiVersion = "6.0.0"
)

// getDataFilePath returns the path to the data file based on the current date
func getDataFilePath() string {
	currentDate := time.Now().Format("20060102")
	return fmt.Sprintf("%s_%s.xml", currentDate, apiVersion)
}

// needsDownload checks if the XML file should be downloaded
// File is downloaded when it doesn't exist or is from a previous day
func needsDownload(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		log.Printf("Error checking file: %v", err)
		return true
	}

	expected := getDataFilePath()
	return filepath.Base(filePath) != filepath.Base(expected)
}

// cleanupAllFiles removes all XML files from the data directory except the current day's file
func cleanupAllFiles(dataDir string) error {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("error reading data directory: %w", err)
	}

	currentFile := getDataFilePath()
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xml") {
			// Skip the current day's file
			if filepath.Base(file.Name()) == filepath.Base(currentFile) {
				continue
			}
			
			filePath := filepath.Join(dataDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error deleting file %s: %v", filePath, err)
			} else {
				log.Printf("Deleted file: %s", filePath)
			}
		}
	}
	return nil
}

// downloadXMLFile downloads the XML file from the medicinal products registry
func downloadXMLFile(url, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	log.Printf("Downloading file from %s...", url)
	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid HTTP status: %d", resp.StatusCode)
	}

	tempFile := filePath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %w", err)
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("error saving data: %w", err)
	}

	out.Close()

	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("cannot rename file: %w", err)
	}

	log.Printf("Downloaded %d bytes in %v", n, time.Since(startTime))

	// Clean up old files only after successful download
	if err := cleanupAllFiles(dir); err != nil {
		log.Printf("Warning: error during cleanup: %v", err)
	}

	return nil
}

// ensureDataFile ensures that the XML file is available and up to date
func ensureDataFile(providedFile string) (string, error) {
	if providedFile != "" && providedFile != getDataFilePath() {
		if _, err := os.Stat(providedFile); err == nil {
			log.Printf("Using user-provided file: %s", providedFile)
			return providedFile, nil
		}
		log.Printf("Provided file does not exist: %s", providedFile)
	}

	filePath := getDataFilePath()

	if needsDownload(filePath) {
		log.Printf("File %s needs to be downloaded", filePath)
		if err := downloadXMLFile(xmlURL, filePath); err != nil {
			return "", fmt.Errorf("cannot download file: %w", err)
		}
	} else {
		log.Printf("Using existing file: %s", filePath)
	}

	return filePath, nil
}

func main() {
	xmlFileFlag := flag.String("file", "", "Optional path to XML file with medicinal products data")
	port := flag.String("port", "1532", "Port to run the HTTP server on")
	flag.Parse()

	xmlFile, err := ensureDataFile(*xmlFileFlag)
	if err != nil {
		log.Fatalf("Error preparing data file: %v", err)
	}

	db := database.NewProductDatabase()
	startTime := time.Now()

	log.Printf("Loading products from %s...", xmlFile)
	if err := db.LoadFromFile(xmlFile); err != nil {
		log.Fatalf("Error loading products: %v", err)
	}

	stats := db.GetStatistics()
	log.Printf("Loaded %d products in %v", stats["liczbaProdukow"], time.Since(startTime))

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.LoadHTMLGlob("templates/*")

	staticDir := "static"
	if _, err := os.Stat(staticDir); !os.IsNotExist(err) {
		router.Static("/static", staticDir)
	}

	handler := api.NewHandler(db)
	handler.RegisterRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	log.Printf("Starting HTTP server on port %s...", *port)
	if err := router.Run(":" + *port); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
