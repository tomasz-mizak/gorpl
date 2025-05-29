// Package api contains HTTP handlers for the API
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorpl/internal/model"
)

// GetUnitboxProductByGtin handles requests for product info in unitbox format by GTIN/EAN
func (h *Handler) GetUnitboxProductByGtin(c *gin.Context) {
	// Get the GTIN from the URL query parameters
	gtin := c.Query("gtin")
	if gtin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing GTIN parameter"})
		return
	}

	// Find the product
	productInfo := h.DB.FindByGtin(gtin)
	if productInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Convert to MedicationTypeRplDto format
	rplProduct := model.ConvertToMedicationTypeRplDto(productInfo)

	c.JSON(http.StatusOK, rplProduct)
}

// SearchUnitboxProductsByName handles search requests for products in unitbox format by name
func (h *Handler) SearchUnitboxProductsByName(c *gin.Context) {
	// Get the query from the URL query parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
		return
	}

	// Search for products
	results := h.DB.SearchByName(query)

	// Convert each result to MedicationTypeRplDto format
	var rplProducts []*model.MedicationTypeRplDto
	for _, product := range results {
		rplProduct := model.ConvertToMedicationTypeRplDto(product)
		if rplProduct != nil {
			rplProducts = append(rplProducts, rplProduct)
		}
	}

	c.JSON(http.StatusOK, rplProducts)
}

// GetSimplifiedMedications handles requests for simplified medication format
func (h *Handler) GetSimplifiedMedications(c *gin.Context) {
	// Get the query from the URL query parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
		return
	}

	// Search for products by name
	resultsByName := h.DB.SearchByName(query)

	// Search for products by GTIN
	resultsByGtin := []*model.ProductInfo{}
	if product := h.DB.FindByGtin(query); product != nil {
		resultsByGtin = append(resultsByGtin, product)
	}

	// Combine and deduplicate results
	seenProducts := make(map[model.BigIntAsString]bool)
	var allResults []*model.ProductInfo

	// Add results from name search
	for _, product := range resultsByName {
		if !seenProducts[product.Product.ID] {
			seenProducts[product.Product.ID] = true
			allResults = append(allResults, product)
		}
	}

	// Add results from GTIN search
	for _, product := range resultsByGtin {
		if !seenProducts[product.Product.ID] {
			seenProducts[product.Product.ID] = true
			allResults = append(allResults, product)
		}
	}

	// Convert results to simplified format
	var simplifiedResults []model.SimplifiedMedicationDto
	for _, product := range allResults {
		if product.Product != nil && product.Package != nil && product.Package.KodGTIN != "" {
			simplifiedResults = append(simplifiedResults, model.SimplifiedMedicationDto{
				TradeName: string(product.Product.NazwaProduktu),
				EanCode:   string(product.Package.KodGTIN),
			})
		}
	}

	c.JSON(http.StatusOK, simplifiedResults)
}

// GetAllSimplifiedMedications handles requests for all medications in simplified format
func (h *Handler) GetAllSimplifiedMedications(c *gin.Context) {
	// Get all products from the database
	results := h.DB.GetAllProducts()

	// Convert results to simplified format
	var simplifiedResults []model.SimplifiedMedicationDto
	for _, product := range results {
		if product.Product != nil && product.Package != nil && product.Package.KodGTIN != "" {
			simplifiedResults = append(simplifiedResults, model.SimplifiedMedicationDto{
				TradeName: string(product.Product.NazwaProduktu),
				EanCode:   string(product.Package.KodGTIN),
			})
		}
	}

	c.JSON(http.StatusOK, simplifiedResults)
}

// RegisterUnitboxRoutes registers all Unitbox API routes
func (h *Handler) RegisterUnitboxRoutes(router *gin.Engine) {
	// API v1 Group for UnitBox
	apiV1 := router.Group("/api/v1/unitbox")
	{
		apiV1.GET("/product", h.GetUnitboxProductByGtin)
		apiV1.GET("/search", h.SearchUnitboxProductsByName)
		apiV1.GET("/simplified", h.GetSimplifiedMedications)
		apiV1.GET("/simplified/all", h.GetAllSimplifiedMedications)
	}
}
