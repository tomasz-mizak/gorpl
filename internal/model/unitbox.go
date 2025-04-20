// Package model contains data structures for representing medicinal products
package model

import (
	"fmt"
	"strings"
)

// MedicationTypeRplDto represents a medicinal product in the format compatible with UnitBox RPL service
type MedicationTypeRplDto struct {
	TradeName         string `json:"tradeName,omitempty"`
	InternationalName string `json:"internationalName,omitempty"`
	Form              string `json:"form,omitempty"`
	Strength          string `json:"strength,omitempty"`
	Unit              string `json:"unit,omitempty"`
	StrengthUnit      string `json:"strengthUnit,omitempty"`
	Manufacturer      string `json:"manufacturer,omitempty"`
	EanCode           string `json:"eanCode,omitempty"`
	AtcCode           string `json:"atcCode,omitempty"`
	Amount            int    `json:"amount,omitempty"`
	AmountUnit        string `json:"amountUnit,omitempty"`
}

// ConvertToMedicationTypeRplDto converts ProductInfo to MedicationTypeRplDto
func ConvertToMedicationTypeRplDto(product *ProductInfo) *MedicationTypeRplDto {
	if product == nil || product.Product == nil || product.Package == nil {
		return nil
	}

	// Extract ATC code from the product (if available)
	var atcCode string
	if product.Product.KodyATC != nil && len(product.Product.KodyATC.KodATC) > 0 {
		atcCode = string(product.Product.KodyATC.KodATC[0])
	}

	// Extract package type and amount
	var amount int = 1 // Default value
	var amountUnit string
	if product.Package.JednostkiOpakowania != nil && len(product.Package.JednostkiOpakowania.JednostkaOpakowania) > 0 {
		unit := product.Package.JednostkiOpakowania.JednostkaOpakowania[0]

		// Try to convert amount to int
		if unit.LiczbaOpakowan != "" {
			fmt.Sscanf(string(unit.LiczbaOpakowan), "%d", &amount)
		}

		// Get amount unit (if available)
		amountUnit = string(unit.JednostkaPojemnosci)
	}

	// Extract manufacturer name
	var manufacturer string
	if product.Product.PodmiotOdpowiedzialny != "" {
		manufacturer = product.Product.PodmiotOdpowiedzialny
	} else if product.Product.DaneOWytworcy != nil && len(product.Product.DaneOWytworcy.Wytworcy) > 0 {
		manufacturer = product.Product.DaneOWytworcy.Wytworcy[0].NazwaWytworcyImportera
	}

	// Parse strength and unit from the Moc field
	strength, unit := parseStrengthUnit(product.Product.Moc)

	return &MedicationTypeRplDto{
		TradeName:         string(product.Product.NazwaProduktu),
		InternationalName: string(product.Product.NazwaPowszechnieStosowana),
		Form:              string(product.Product.NazwaPostaciFarmaceutycznej),
		Strength:          strength,
		Unit:              unit,
		StrengthUnit:      product.Product.Moc,
		Manufacturer:      manufacturer,
		EanCode:           string(product.Package.KodGTIN),
		AtcCode:           atcCode,
		Amount:            amount,
		AmountUnit:        amountUnit,
	}
}

// parseStrengthUnit attempts to parse strength and unit from a combined string
// For example, "10 mg" would return "10", "mg"
func parseStrengthUnit(combined string) (string, string) {
	if combined == "" {
		return "", ""
	}

	// Common patterns:
	// "10 mg"
	// "10 mg + 5 mg"
	// "10 mg/ml"

	// Simple case - try to split on first space
	parts := strings.SplitN(combined, " ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	// Try to find numeric part
	var strength, unit string
	for i, c := range combined {
		if (c < '0' || c > '9') && c != '.' && c != ',' {
			// First non-digit character
			strength = combined[:i]
			unit = combined[i:]
			break
		}
	}

	if strength == "" {
		// No numeric part found
		return "", combined
	}

	return strings.TrimSpace(strength), strings.TrimSpace(unit)
}
