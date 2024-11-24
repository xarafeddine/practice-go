package main

import (
	"fmt"
	"math"
)

// Types and constants
type CoffeeSize string
type CoffeeType string
type AddOnType string

const (
	// Coffee sizes
	SizeSmall  CoffeeSize = "SMALL"
	SizeMedium CoffeeSize = "MEDIUM"
	SizeLarge  CoffeeSize = "LARGE"

	// Coffee types
	TypeEspresso   CoffeeType = "ESPRESSO"
	TypeLatte      CoffeeType = "LATTE"
	TypeCappuccino CoffeeType = "CAPPUCCINO"
	TypeAmericano  CoffeeType = "AMERICANO"

	// Add-on types
	AddOnExtraShot    AddOnType = "EXTRA_SHOT"
	AddOnWhippedCream AddOnType = "WHIPPED_CREAM"
	AddOnCaramel      AddOnType = "CARAMEL"
	AddOnChocolate    AddOnType = "CHOCOLATE"
	AddOnSoyMilk      AddOnType = "SOY_MILK"
)

// Structs
type Coffee struct {
	Type CoffeeType
	Size CoffeeSize
}

type AddOn struct {
	Type  AddOnType
	Price float64
}

type Order struct {
	Coffee     Coffee
	AddOns     []AddOn
	Promotions []Promotion
}

type Promotion interface {
	Name() string
	Description() string
	IsActive() bool
	Apply(price float64, addOns []AddOn) float64
}

// Price configurations
var coffeeBasePrices = map[CoffeeType]map[CoffeeSize]float64{
	TypeEspresso: {
		SizeSmall:  2.50,
		SizeMedium: 3.00,
		SizeLarge:  3.50,
	},
	TypeLatte: {
		SizeSmall:  3.00,
		SizeMedium: 3.50,
		SizeLarge:  4.00,
	},
	TypeCappuccino: {
		SizeSmall:  3.00,
		SizeMedium: 3.50,
		SizeLarge:  4.00,
	},
	TypeAmericano: {
		SizeSmall:  2.00,
		SizeMedium: 2.50,
		SizeLarge:  3.00,
	},
}

var addOns = map[AddOnType]AddOn{
	AddOnExtraShot:    {Type: AddOnExtraShot, Price: 0.50},
	AddOnWhippedCream: {Type: AddOnWhippedCream, Price: 0.75},
	AddOnCaramel:      {Type: AddOnCaramel, Price: 0.50},
	AddOnChocolate:    {Type: AddOnChocolate, Price: 0.50},
	AddOnSoyMilk:      {Type: AddOnSoyMilk, Price: 1.00},
}

// Promotion implementations
type PercentageDiscount struct {
	name        string
	description string
	active      bool
	percentage  float64
}

func NewPercentageDiscount(name, description string, active bool, percentage float64) *PercentageDiscount {
	return &PercentageDiscount{
		name:        name,
		description: description,
		active:      active,
		percentage:  percentage,
	}
}

func (p *PercentageDiscount) Name() string        { return p.name }
func (p *PercentageDiscount) Description() string { return p.description }
func (p *PercentageDiscount) IsActive() bool      { return p.active }
func (p *PercentageDiscount) Apply(price float64, _ []AddOn) float64 {
	return price * (1 - p.percentage/100)
}

type FreeExpensiveAddOn struct {
	name        string
	description string
	active      bool
}

func NewFreeExpensiveAddOn(name, description string, active bool) *FreeExpensiveAddOn {
	return &FreeExpensiveAddOn{
		name:        name,
		description: description,
		active:      active,
	}
}

func (f *FreeExpensiveAddOn) Name() string        { return f.name }
func (f *FreeExpensiveAddOn) Description() string { return f.description }
func (f *FreeExpensiveAddOn) IsActive() bool      { return f.active }
func (f *FreeExpensiveAddOn) Apply(price float64, addOns []AddOn) float64 {
	if len(addOns) == 0 {
		return price
	}

	mostExpensive := addOns[0]
	for _, addOn := range addOns {
		if addOn.Price > mostExpensive.Price {
			mostExpensive = addOn
		}
	}

	return price - mostExpensive.Price
}

// OrderBuilder provides a fluent interface for building orders
type OrderBuilder struct {
	order Order
	err   error
}

func NewOrder(coffeeType CoffeeType, size CoffeeSize) *OrderBuilder {
	return &OrderBuilder{
		order: Order{
			Coffee: Coffee{
				Type: coffeeType,
				Size: size,
			},
		},
	}
}

func (b *OrderBuilder) AddAddOn(addOnType AddOnType) *OrderBuilder {
	if b.err != nil {
		return b
	}

	addOn, exists := addOns[addOnType]
	if !exists {
		b.err = fmt.Errorf("invalid add-on type: %s", addOnType)
		return b
	}

	b.order.AddOns = append(b.order.AddOns, addOn)
	return b
}

func (b *OrderBuilder) AddPromotion(promotion Promotion) *OrderBuilder {
	if b.err != nil {
		return b
	}

	b.order.Promotions = append(b.order.Promotions, promotion)
	return b
}

func (b *OrderBuilder) Build() (*Order, error) {
	if b.err != nil {
		return nil, b.err
	}
	return &b.order, nil
}

// OrderSummary represents the final order details
type OrderSummary struct {
	Coffee     Coffee
	AddOns     []AddOn
	Promotions []Promotion
	FinalPrice float64
}

// CalculatePrice calculates the final price of an order
func CalculatePrice(order *Order) (float64, error) {
	basePrice, exists := coffeeBasePrices[order.Coffee.Type][order.Coffee.Size]
	if !exists {
		return 0, fmt.Errorf("invalid coffee type or size")
	}

	// Add add-ons
	for _, addOn := range order.AddOns {
		basePrice += addOn.Price
	}

	// Apply promotions
	finalPrice := basePrice
	for _, promotion := range order.Promotions {
		if promotion.IsActive() {
			finalPrice = promotion.Apply(finalPrice, order.AddOns)
		}
	}

	return math.Round(finalPrice*100) / 100, nil
}

// GetOrderSummary generates a summary of the order
func GetOrderSummary(order *Order) (*OrderSummary, error) {
	finalPrice, err := CalculatePrice(order)
	if err != nil {
		return nil, err
	}

	return &OrderSummary{
		Coffee:     order.Coffee,
		AddOns:     order.AddOns,
		Promotions: order.Promotions,
		FinalPrice: finalPrice,
	}, nil
}

// Example usage
func CoffeeOrderApp() {
	// Create a new order
	builder := NewOrder(TypeLatte, SizeMedium).
		AddAddOn(AddOnExtraShot).
		AddAddOn(AddOnWhippedCream).
		AddPromotion(NewPercentageDiscount("20% Off", "Get 20% off your order", true, 20))

	order, err := builder.Build()
	if err != nil {
		fmt.Printf("Error building order: %v\n", err)
		return
	}

	summary, err := GetOrderSummary(order)
	if err != nil {
		fmt.Printf("Error getting summary: %v\n", err)
		return
	}

	fmt.Printf("Order Summary:\n")
	fmt.Printf("Coffee: %s %s\n", summary.Coffee.Size, summary.Coffee.Type)
	fmt.Printf("Add-ons: %d\n", len(summary.AddOns))
	fmt.Printf("Final Price: $%.2f\n", summary.FinalPrice)
}
