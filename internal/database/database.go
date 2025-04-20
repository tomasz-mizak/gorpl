// Package database zawiera kod obsługi bazy danych produktów leczniczych
package database

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"gorpl/internal/model"
)

// ProductRepository defines interface for product database operations
type ProductRepository interface {
	LoadFromFile(filename string) error
	FindByGtin(gtin string) *model.ProductInfo
	SearchByName(query string) []*model.ProductInfo
	GetStatistics() map[string]interface{}
}

// ProductDatabase holds the database of medical products and provides methods to search it
type ProductDatabase struct {
	produkty *model.ProduktyLecznicze
	// Map for quick lookups by GTIN (EAN)
	gtinIndex map[string]*model.ProductInfo
	mutex     sync.RWMutex
}

// Make sure ProductDatabase implements ProductRepository
var _ ProductRepository = (*ProductDatabase)(nil)

// NewProductDatabase creates a new product database
func NewProductDatabase() *ProductDatabase {
	return &ProductDatabase{
		gtinIndex: make(map[string]*model.ProductInfo),
	}
}

// LoadFromFile loads the products from an XML file
func (db *ProductDatabase) LoadFromFile(filename string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a decoder
	decoder := xml.NewDecoder(file)

	// Create an instance of ProduktyLecznicze to store the parsed data
	var produkty model.ProduktyLecznicze

	// Parse the XML
	err = decoder.Decode(&produkty)
	if err != nil {
		return fmt.Errorf("error decoding XML: %w", err)
	}

	db.produkty = &produkty

	// Build the GTIN index
	db.buildGtinIndex()

	return nil
}

// buildGtinIndex creates an index of products by GTIN for fast lookups
func (db *ProductDatabase) buildGtinIndex() {
	// Clear existing index
	db.gtinIndex = make(map[string]*model.ProductInfo)

	// Build new index
	for i := range db.produkty.ProduktyLecznicze {
		product := &db.produkty.ProduktyLecznicze[i]

		if product.Opakowania == nil {
			continue
		}

		for j := range product.Opakowania.Opakowanie {
			pkg := &product.Opakowania.Opakowanie[j]

			// Skip deleted packages
			if pkg.Skasowane == "TAK" {
				continue
			}

			// Index by GTIN/EAN
			if pkg.KodGTIN != "" {
				db.gtinIndex[string(pkg.KodGTIN)] = &model.ProductInfo{
					Product: product,
					Package: pkg,
				}
			}

			// Also index by foreign GTINs if available
			if pkg.ZgodyPrezesa != nil {
				for _, zgoda := range pkg.ZgodyPrezesa.ZgodaPrezesa {
					if zgoda.GTINZagraniczne != nil {
						for _, gtin := range zgoda.GTINZagraniczne.GTINZagraniczny {
							if gtin.Numer != "" {
								db.gtinIndex[gtin.Numer] = &model.ProductInfo{
									Product: product,
									Package: pkg,
								}
							}
						}
					}
				}
			}
		}
	}

	log.Printf("Built GTIN index with %d entries", len(db.gtinIndex))
}

// FindByGtin finds a product by its GTIN/EAN code
func (db *ProductDatabase) FindByGtin(gtin string) *model.ProductInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return db.gtinIndex[gtin]
}

// GetStatistics returns statistics about the database
func (db *ProductDatabase) GetStatistics() map[string]interface{} {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return map[string]interface{}{
		"stanNaDzien":       db.produkty.StanNaDzien,
		"liczbaProdukow":    len(db.produkty.ProduktyLecznicze),
		"liczbaIndeksowEAN": len(db.gtinIndex),
	}
}

// SearchByName searches for products by name (partial match)
func (db *ProductDatabase) SearchByName(query string) []*model.ProductInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// If query is empty, return empty result
	if query == "" {
		return nil
	}

	var results []*model.ProductInfo
	seenProducts := make(map[model.BigIntAsString]bool)

	// Search through all products
	for i := range db.produkty.ProduktyLecznicze {
		product := &db.produkty.ProduktyLecznicze[i]

		// Skip if we've already seen this product
		if seenProducts[product.ID] {
			continue
		}

		// Check if product name contains the query (case-insensitive)
		if containsIgnoreCase(string(product.NazwaProduktu), query) ||
			containsIgnoreCase(string(product.NazwaPowszechnieStosowana), query) {

			seenProducts[product.ID] = true

			// Find first non-deleted package for this product
			if product.Opakowania != nil {
				for j := range product.Opakowania.Opakowanie {
					pkg := &product.Opakowania.Opakowanie[j]

					if pkg.Skasowane != "TAK" {
						results = append(results, &model.ProductInfo{
							Product: product,
							Package: pkg,
						})
						break
					}
				}
			}
		}
	}

	return results
}

// Helper function to check if a string contains another string (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
