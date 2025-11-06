package valueobjects

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Money represents a monetary value
type Money struct {
	Amount   int64
	Currency string
}

// NewMoney creates a new Money instance
func NewMoney(amount int64) *Money {
	return &Money{
		Amount:   amount,
		Currency: "IDR", // Indonesian Rupiah as default
	}
}

// NewMoneyWithCurrency creates a new Money instance with specified currency
func NewMoneyWithCurrency(amount int64, currency string) (*Money, error) {
	if currency == "" {
		return nil, errors.New("currency is required")
	}
	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Add adds another Money value to this one
func (m *Money) Add(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, fmt.Errorf("cannot add different currencies: %s and %s", m.Currency, other.Currency)
	}

	return &Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract subtracts another Money value from this one
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, fmt.Errorf("cannot subtract different currencies: %s and %s", m.Currency, other.Currency)
	}

	result := m.Amount - other.Amount
	if result < 0 {
		return nil, errors.New("resulting amount cannot be negative")
	}

	return &Money{
		Amount:   result,
		Currency: m.Currency,
	}, nil
}

// Multiply multiplies the money value by a factor
func (m *Money) Multiply(factor float64) (*Money, error) {
	if factor < 0 {
		return nil, errors.New("multiplication factor cannot be negative")
	}

	result := int64(float64(m.Amount) * factor)
	return &Money{
		Amount:   result,
		Currency: m.Currency,
	}, nil
}

// IsZero checks if the amount is zero
func (m *Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive checks if the amount is positive
func (m *Money) IsPositive() bool {
	return m.Amount > 0
}

// IsNegative checks if the amount is negative
func (m *Money) IsNegative() bool {
	return m.Amount < 0
}

// Equals checks if two Money values are equal
func (m *Money) Equals(other *Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

// GreaterThan checks if this Money is greater than another
func (m *Money) GreaterThan(other *Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount > other.Amount
}

// GreaterThanOrEqual checks if this Money is greater than or equal to another
func (m *Money) GreaterThanOrEqual(other *Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount >= other.Amount
}

// LessThan checks if this Money is less than another
func (m *Money) LessThan(other *Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount < other.Amount
}

// LessThanOrEqual checks if this Money is less than or equal to another
func (m *Money) LessThanOrEqual(other *Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount <= other.Amount
}

// String returns a formatted string representation
func (m *Money) String() string {
	return fmt.Sprintf("%s %d", m.Currency, m.Amount)
}

// FormatIndonesian formats the money for Indonesian Rupiah display
func (m *Money) FormatIndonesian() string {
	if m.Currency != "IDR" {
		return m.String()
	}

	// Format with thousand separators
	str := strconv.FormatInt(m.Amount, 10)
	var result []rune

	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, digit)
	}

	return fmt.Sprintf("Rp %s", string(result))
}

// ParseMoney parses a string into Money
func ParseMoney(amountStr string, currency string) (*Money, error) {
	if amountStr == "" {
		return nil, errors.New("amount string is empty")
	}

	// Remove common formatting characters
	cleaned := strings.ReplaceAll(amountStr, ",", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.TrimSpace(cleaned)

	amount, err := strconv.ParseInt(cleaned, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	return NewMoneyWithCurrency(amount, currency)
}

// MoneyValidationRules represents validation rules for Money
type MoneyValidationRules struct {
	MinAmount   int64
	MaxAmount   int64
	Currencies  []string
	AllowZero   bool
}

// Validate validates the Money against given rules
func (m *Money) Validate(rules MoneyValidationRules) error {
	// Check currency
	if len(rules.Currencies) > 0 {
		found := false
		for _, currency := range rules.Currencies {
			if m.Currency == currency {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid currency: %s", m.Currency)
		}
	}

	// Check amount range
	if m.Amount < rules.MinAmount {
		return fmt.Errorf("amount %d is below minimum %d", m.Amount, rules.MinAmount)
	}

	if rules.MaxAmount > 0 && m.Amount > rules.MaxAmount {
		return fmt.Errorf("amount %d is above maximum %d", m.Amount, rules.MaxAmount)
	}

	// Check zero allowance
	if !rules.AllowZero && m.Amount == 0 {
		return errors.New("zero amount is not allowed")
	}

	return nil
}

// Common validation rules
var (
	// IDRValidationRules for Indonesian Rupiah
	IDRValidationRules = MoneyValidationRules{
		MinAmount:   0,
		MaxAmount:   0, // No maximum
		Currencies:  []string{"IDR"},
		AllowZero:   true,
	}

	// PositiveIDRValidationRules for positive amounts only
	PositiveIDRValidationRules = MoneyValidationRules{
		MinAmount:   1,
		MaxAmount:   0, // No maximum
		Currencies:  []string{"IDR"},
		AllowZero:   false,
	}

	// TicketCostValidationRules for ticket costs
	TicketCostValidationRules = MoneyValidationRules{
		MinAmount:   0,
		MaxAmount:   10000000000, // 10 billion IDR
		Currencies:  []string{"IDR"},
		AllowZero:   true,
	}
)