// Package api contains handlers for the API
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorpl/internal/database"
	"gorpl/internal/model"
)

// Handler structure holds dependencies for API handlers
type Handler struct {
	DB database.ProductRepository
}

// NewHandler creates a new Handler instance
func NewHandler(db database.ProductRepository) *Handler {
	return &Handler{DB: db}
}

// GetProductByGtin handles requests for product info by GTIN/EAN
func (h *Handler) GetProductByGtin(c *gin.Context) {
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

	// Create a response structure to ensure proper JSON serialization
	response := struct {
		Product *model.ProduktLeczniczy `json:"product"`
		Package *model.Opakowanie       `json:"package"`
	}{
		Product: productInfo.Product,
		Package: productInfo.Package,
	}

	c.JSON(http.StatusOK, response)
}

// SearchProductsByName handles search requests by product name
func (h *Handler) SearchProductsByName(c *gin.Context) {
	// Get the query from the URL query parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
		return
	}

	// Search for products
	results := h.DB.SearchByName(query)

	// Transform results to ensure proper JSON serialization
	var response []struct {
		Product *model.ProduktLeczniczy `json:"product"`
		Package *model.Opakowanie       `json:"package"`
	}

	for _, result := range results {
		response = append(response, struct {
			Product *model.ProduktLeczniczy `json:"product"`
			Package *model.Opakowanie       `json:"package"`
		}{
			Product: result.Product,
			Package: result.Package,
		})
	}

	// Return results (even if empty)
	c.JSON(http.StatusOK, response)
}

// GetStats handles requests for database statistics
func (h *Handler) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, h.DB.GetStatistics())
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// API Group
	api := router.Group("/api/v1")
	{
		api.GET("/product", h.GetProductByGtin)
		api.GET("/search", h.SearchProductsByName)
		api.GET("/stats", h.GetStats)
	}

	// Register Unitbox specific routes
	h.RegisterUnitboxRoutes(router)
}
