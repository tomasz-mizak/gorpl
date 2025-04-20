// Package main zawiera główny kod serwera HTTP
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorpl/internal/database"
)

func main() {
	// Define command-line flags
	xmlFile := flag.String("file", "20250419_6.0.0.xml", "Path to the XML file with medical products data")
	port := flag.String("port", "8080", "Port to run the HTTP server on")
	flag.Parse()
	
	// Create a new product database
	db := database.NewProductDatabase()
	
	// Start time for performance measurement
	startTime := time.Now()
	
	// Check if the file exists
	if _, err := os.Stat(*xmlFile); os.IsNotExist(err) {
		log.Fatalf("Error: File %s does not exist", *xmlFile)
	}
	
	// Load the products from the XML file
	log.Printf("Loading products from %s...", *xmlFile)
	if err := db.LoadFromFile(*xmlFile); err != nil {
		log.Fatalf("Error loading products: %v", err)
	}
	
	// Get statistics to display product count
	stats := db.GetStatistics()
	log.Printf("Loaded %d products in %v", stats["liczbaProdukow"], time.Since(startTime))
	
	// Set up HTTP routes
	http.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		handleProductByGtin(w, r, db)
	})
	
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		handleSearchByName(w, r, db)
	})
	
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		handleStats(w, r, db)
	})
	
	// Simple home page
	http.HandleFunc("/", handleHomePage)
	
	// Start the HTTP server
	log.Printf("Starting HTTP server on port %s...", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

// handleProductByGtin handles requests for product info by GTIN/EAN
func handleProductByGtin(w http.ResponseWriter, r *http.Request, db *database.ProductDatabase) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get the GTIN from the URL query parameters
	gtin := r.URL.Query().Get("gtin")
	if gtin == "" {
		http.Error(w, "Missing GTIN parameter", http.StatusBadRequest)
		return
	}
	
	// Find the product
	productInfo := db.FindByGtin(gtin)
	if productInfo == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	
	// Encode the product info as JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err := encoder.Encode(productInfo); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleSearchByName handles search requests by product name
func handleSearchByName(w http.ResponseWriter, r *http.Request, db *database.ProductDatabase) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get the query from the URL query parameters
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}
	
	// Search for products
	results := db.SearchByName(query)
	if len(results) == 0 {
		// Return empty array instead of 404
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "[]")
		return
	}
	
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	
	// Encode the results as JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err := encoder.Encode(results); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleStats handles requests for database statistics
func handleStats(w http.ResponseWriter, r *http.Request, db *database.ProductDatabase) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get statistics
	stats := db.GetStatistics()
	
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	
	// Encode the stats as JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err := encoder.Encode(stats); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleHomePage serves the home page HTML
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Baza Produktów Leczniczych</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        h1 { color: #333; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
        input[type="text"] { padding: 8px; width: 300px; }
        button { padding: 8px 15px; background: #4CAF50; color: white; border: none; cursor: pointer; }
        pre { background: #f4f4f4; padding: 15px; overflow: auto; }
        .tabs { display: flex; margin-bottom: 20px; }
        .tab { padding: 10px 15px; cursor: pointer; background: #f1f1f1; margin-right: 5px; }
        .tab.active { background: #4CAF50; color: white; }
        .tab-content { display: none; }
        .tab-content.active { display: block; }
    </style>
</head>
<body>
    <h1>Baza Produktów Leczniczych</h1>
    
    <div class="tabs">
        <div class="tab active" onclick="showTab('search-gtin')">Wyszukaj po EAN/GTIN</div>
        <div class="tab" onclick="showTab('search-name')">Wyszukaj po nazwie</div>
    </div>
    
    <div id="search-gtin" class="tab-content active">
        <div class="form-group">
            <label for="gtin">Kod EAN/GTIN:</label>
            <input type="text" id="gtin" placeholder="Wprowadź kod EAN/GTIN">
            <button onclick="searchByGtin()">Szukaj</button>
        </div>
    </div>
    
    <div id="search-name" class="tab-content">
        <div class="form-group">
            <label for="product-name">Nazwa produktu:</label>
            <input type="text" id="product-name" placeholder="Wprowadź nazwę produktu">
            <button onclick="searchByName()">Szukaj</button>
        </div>
    </div>
    
    <div id="result">
        <pre id="json-result"></pre>
    </div>

    <script>
        function showTab(tabId) {
            // Hide all tabs
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.remove('active');
            });
            document.querySelectorAll('.tab').forEach(tab => {
                tab.classList.remove('active');
            });
            
            // Show selected tab
            document.getElementById(tabId).classList.add('active');
            document.querySelector(`.tab[onclick="showTab('${tabId}')"]`).classList.add('active');
        }
        
        async function searchByGtin() {
            const gtin = document.getElementById('gtin').value;
            if (!gtin) {
                alert('Wprowadź kod EAN/GTIN');
                return;
            }
            
            try {
                const response = await fetch('/api/product?gtin=' + encodeURIComponent(gtin));
                if (!response.ok) {
                    throw new Error('Nie znaleziono produktu lub wystąpił błąd');
                }
                
                const data = await response.json();
                document.getElementById('json-result').textContent = JSON.stringify(data, null, 2);
            } catch (error) {
                document.getElementById('json-result').textContent = error.message;
            }
        }
        
        async function searchByName() {
            const name = document.getElementById('product-name').value;
            if (!name) {
                alert('Wprowadź nazwę produktu');
                return;
            }
            
            try {
                const response = await fetch('/api/search?query=' + encodeURIComponent(name));
                if (!response.ok) {
                    throw new Error('Wystąpił błąd podczas wyszukiwania');
                }
                
                const data = await response.json();
                if (data.length === 0) {
                    document.getElementById('json-result').textContent = 'Nie znaleziono produktów pasujących do zapytania';
                } else {
                    document.getElementById('json-result').textContent = JSON.stringify(data, null, 2);
                }
            } catch (error) {
                document.getElementById('json-result').textContent = error.message;
            }
        }
    </script>
</body>
</html>
	`)
}