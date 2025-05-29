// Package database contains code for handling medicinal products database
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
	SearchByGtin(gtin string) []*model.ProductInfo
	GetStatistics() map[string]interface{}
	GetAllProducts() []*model.ProductInfo
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

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var produkty model.ProduktyLecznicze

	err = decoder.Decode(&produkty)
	if err != nil {
		return fmt.Errorf("error decoding XML: %w", err)
	}

	db.produkty = &produkty
	db.buildGtinIndex()

	return nil
}

// buildGtinIndex creates an index of products by GTIN for fast lookups
func (db *ProductDatabase) buildGtinIndex() {
	db.gtinIndex = make(map[string]*model.ProductInfo)

	for i := range db.produkty.ProduktyLecznicze {
		product := &db.produkty.ProduktyLecznicze[i]

		if product.Opakowania == nil {
			continue
		}

		for j := range product.Opakowania.Opakowanie {
			pkg := &product.Opakowania.Opakowanie[j]

			if pkg.Skasowane == "TAK" {
				continue
			}

			if pkg.KodGTIN != "" {
				db.gtinIndex[string(pkg.KodGTIN)] = &model.ProductInfo{
					Product: product,
					Package: pkg,
				}
			}

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

	if query == "" {
		return nil
	}

	var results []*model.ProductInfo
	seenProducts := make(map[model.BigIntAsString]bool)

	for i := range db.produkty.ProduktyLecznicze {
		product := &db.produkty.ProduktyLecznicze[i]

		if seenProducts[product.ID] {
			continue
		}

		if containsIgnoreCase(string(product.NazwaProduktu), query) ||
			containsIgnoreCase(string(product.NazwaPowszechnieStosowana), query) {

			seenProducts[product.ID] = true

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

// containsIgnoreCase checks if a string contains another string (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)

	// Split both strings into words
	sWords := strings.Fields(s)
	substrWords := strings.Fields(substr)

	// If substr is a single word, check if it's a prefix of any word in s
	if len(substrWords) == 1 {
		for _, word := range sWords {
			if strings.HasPrefix(word, substr) {
				return true
			}
		}
		return false
	}

	// For multi-word searches, check if all words are present in order
	return strings.Contains(s, substr)
}

// GetAllProducts returns all products from the database
func (db *ProductDatabase) GetAllProducts() []*model.ProductInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	var results []*model.ProductInfo
	seenProducts := make(map[model.BigIntAsString]bool)

	for i := range db.produkty.ProduktyLecznicze {
		product := &db.produkty.ProduktyLecznicze[i]

		if seenProducts[product.ID] {
			continue
		}

		seenProducts[product.ID] = true

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

	return results
}

// SearchByGtin searches for products by GTIN/EAN code (partial match)
func (db *ProductDatabase) SearchByGtin(gtin string) []*model.ProductInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if gtin == "" {
		return nil
	}

	var results []*model.ProductInfo
	seenProducts := make(map[model.BigIntAsString]bool)

	// Search through all products
	for i := range db.produkty.ProduktyLecznicze {
		product := &db.produkty.ProduktyLecznicze[i]

		if seenProducts[product.ID] {
			continue
		}

		if product.Opakowania != nil {
			for j := range product.Opakowania.Opakowanie {
				pkg := &product.Opakowania.Opakowanie[j]

				if pkg.Skasowane == "TAK" {
					continue
				}

				// Check main GTIN
				if pkg.KodGTIN != "" && containsIgnoreCase(string(pkg.KodGTIN), gtin) {
					seenProducts[product.ID] = true
					results = append(results, &model.ProductInfo{
						Product: product,
						Package: pkg,
					})
					break
				}

				// Check foreign GTINs
				if pkg.ZgodyPrezesa != nil {
					for _, zgoda := range pkg.ZgodyPrezesa.ZgodaPrezesa {
						if zgoda.GTINZagraniczne != nil {
							for _, foreignGtin := range zgoda.GTINZagraniczne.GTINZagraniczny {
								if foreignGtin.Numer != "" && containsIgnoreCase(foreignGtin.Numer, gtin) {
									seenProducts[product.ID] = true
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
			}
		}
	}

	return results
}
