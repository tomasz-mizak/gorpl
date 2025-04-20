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

// RegisterUnitboxRoutes registers all Unitbox API routes
func (h *Handler) RegisterUnitboxRoutes(router *gin.Engine) {
	// API v1 Group for UnitBox
	apiV1 := router.Group("/api/v1/unitbox")
	{
		apiV1.GET("/product", h.GetUnitboxProductByGtin)
		apiV1.GET("/search", h.SearchUnitboxProductsByName)
	}
}
