// Package main zawiera główny kod serwera HTTP
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"gorpl/internal/api"
	"gorpl/internal/database"
)

const (
	// URL do pobierania pliku XML z rejestru produktów leczniczych
	xmlURL = "https://rejestry.ezdrowie.gov.pl/api/rpl/medicinal-products/public-pl-report/6.0.0/overall.xml"
	// Wersja API
	apiVersion = "6.0.0"
)

// getDataFilePath zwraca ścieżkę do pliku danych na podstawie aktualnej daty
func getDataFilePath() string {
	currentDate := time.Now().Format("20060102")
	return fmt.Sprintf("%s_%s.xml", currentDate, apiVersion)
}

// needsDownload sprawdza, czy plik XML powinien zostać pobrany
// Plik jest pobierany, gdy nie istnieje lub pochodzi z wcześniejszego dnia
func needsDownload(filePath string) bool {
	// Sprawdź, czy plik istnieje
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// Plik nie istnieje
		return true
	}
	if err != nil {
		// Inny błąd
		log.Printf("Błąd podczas sprawdzania pliku: %v", err)
		return true
	}

	// Pobierz datę z nazwy bieżącego pliku
	expected := getDataFilePath()
	return filepath.Base(filePath) != filepath.Base(expected)
}

// downloadXMLFile pobiera plik XML z rejestru produktów leczniczych
func downloadXMLFile(url, filePath string) error {
	// Utwórz katalog, jeśli nie istnieje
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("nie można utworzyć katalogu: %w", err)
	}

	// Pobierz plik
	log.Printf("Pobieranie pliku z %s...", url)
	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("błąd podczas pobierania: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nieprawidłowy status HTTP: %d", resp.StatusCode)
	}

	// Utwórz plik tymczasowy
	tempFile := filePath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("nie można utworzyć pliku tymczasowego: %w", err)
	}
	defer out.Close()

	// Kopiuj dane
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tempFile) // Usuń plik tymczasowy w przypadku błędu
		return fmt.Errorf("błąd podczas zapisywania danych: %w", err)
	}

	// Zamknij plik przed zmianą nazwy
	out.Close()

	// Zmień nazwę pliku tymczasowego na docelową
	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile) // Usuń plik tymczasowy w przypadku błędu
		return fmt.Errorf("nie można zmienić nazwy pliku: %w", err)
	}

	log.Printf("Pobrano %d bajtów w %v", n, time.Since(startTime))
	return nil
}

// ensureDataFile zapewnia, że plik XML jest dostępny i aktualny
func ensureDataFile(providedFile string) (string, error) {
	// Jeśli plik został podany przez użytkownika, użyj go
	if providedFile != "" && providedFile != getDataFilePath() {
		if _, err := os.Stat(providedFile); err == nil {
			log.Printf("Używam pliku podanego przez użytkownika: %s", providedFile)
			return providedFile, nil
		}
		log.Printf("Podany plik nie istnieje: %s", providedFile)
	}

	// Użyj domyślnej ścieżki pliku
	filePath := getDataFilePath()

	// Sprawdź, czy plik wymaga pobrania
	if needsDownload(filePath) {
		log.Printf("Plik %s wymaga pobrania", filePath)
		if err := downloadXMLFile(xmlURL, filePath); err != nil {
			return "", fmt.Errorf("nie można pobrać pliku: %w", err)
		}
	} else {
		log.Printf("Używam istniejącego pliku: %s", filePath)
	}

	return filePath, nil
}

func main() {
	// Define command-line flags
	xmlFileFlag := flag.String("file", "", "Opcjonalna ścieżka do pliku XML z danymi produktów leczniczych")
	port := flag.String("port", "1532", "Port to run the HTTP server on")
	flag.Parse()

	// Zapewnij, że plik XML jest dostępny
	xmlFile, err := ensureDataFile(*xmlFileFlag)
	if err != nil {
		log.Fatalf("Błąd przygotowania pliku danych: %v", err)
	}

	// Create a new product database
	db := database.NewProductDatabase()

	// Start time for performance measurement
	startTime := time.Now()

	// Load the products from the XML file
	log.Printf("Ładowanie produktów z %s...", xmlFile)
	if err := db.LoadFromFile(xmlFile); err != nil {
		log.Fatalf("Błąd podczas ładowania produktów: %v", err)
	}

	// Get statistics to display product count
	stats := db.GetStatistics()
	log.Printf("Załadowano %d produktów w %v", stats["liczbaProdukow"], time.Since(startTime))

	// Initialize Gin
	router := gin.Default()

	// No trust for proxies
	router.SetTrustedProxies(nil)

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Serve static files if they exist
	staticDir := "static"
	if _, err := os.Stat(staticDir); !os.IsNotExist(err) {
		router.Static("/static", staticDir)
	}

	// Create API handler and register routes
	handler := api.NewHandler(db)
	handler.RegisterRoutes(router)

	// Home page route
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Start the HTTP server
	log.Printf("Uruchamianie serwera HTTP na porcie %s...", *port)
	if err := router.Run(":" + *port); err != nil {
		log.Fatalf("Błąd podczas uruchamiania serwera HTTP: %v", err)
	}
}
