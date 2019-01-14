package gorm

import (
	"cafapp-returns/logger"
	"fmt"
	"strings"
)

func properFoodName(in string) string {
	return strings.Title(strings.Replace(in, "-", " ", -1))
}

func initData() {
	allFoods := []Product{}

	// chicken2
	chicken2 := Product{
		SKU:          "C2",
		Name:         fmt.Sprintf("2-chicken"),
		PriceInCents: 590,
		Description:  fmt.Sprintf("2 Chicken Strips, French Fries, Cooler Drink"),
		Status:       ProductStatusOnShelf,
	}
	chicken2.DisplayName = properFoodName("chicken tender meal 2pc")
	chicken2.DescriptionHTML = strings.Replace(chicken2.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, chicken2)

	// chickend4
	chicken4 := Product{
		SKU:          "C4",
		Name:         fmt.Sprintf("4-chicken"),
		PriceInCents: 790,
		Description:  fmt.Sprintf("4 Chicken Strips, French Fries, Cooler Drink"),
		Status:       ProductStatusOnShelf,
	}
	chicken4.DisplayName = properFoodName("chicken tender meal 4pc")
	chicken4.DescriptionHTML = strings.Replace(chicken4.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, chicken4)

	// burger
	burger := Product{
		SKU:          "B0",
		Name:         fmt.Sprintf("burger"),
		PriceInCents: 715,
		Description:  fmt.Sprintf("Cheeseburger w/ lettuce & tomato), French Fries, Cooler Drink"),
		Status:       ProductStatusOnShelf,
	}
	burger.DisplayName = properFoodName("burger meal")
	burger.DescriptionHTML = strings.Replace(burger.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, burger)

	// turkeysandwich
	turkeysandwich := Product{
		SKU:          "S0",
		Name:         fmt.Sprintf("turkey-sandwich"),
		PriceInCents: 835,
		Description:  fmt.Sprintf("Turkey Sub Sandwich, French Fries, Cooler Drink, Toppings include: turkey + lettuce + tomato, with mayo & mustard in packets"),
		Status:       ProductStatusOnShelf,
	}
	turkeysandwich.DisplayName = properFoodName("turkey sandwich meal")
	turkeysandwich.DescriptionHTML = strings.Replace(turkeysandwich.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, turkeysandwich)

	// veggiesandwich
	veggiesandwich := Product{
		SKU:          "S1",
		Name:         fmt.Sprintf("veggie-sandwich"),
		PriceInCents: 835,
		Description:  fmt.Sprintf("Veggie Sub Sandwich, French Fries, Cooler Drink, Toppings include: hummus + cucumber + lettuce + tomato + green pepper, with mayo & mustard in packets"),
		Status:       ProductStatusOnShelf,
	}
	veggiesandwich.DisplayName = properFoodName("veggie sandwich meal")
	veggiesandwich.DescriptionHTML = strings.Replace(veggiesandwich.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, veggiesandwich)

	// peppizza
	peppizza := Product{
		SKU:          "P0",
		Name:         fmt.Sprintf("pep-pizza"),
		PriceInCents: 730,
		Description:  fmt.Sprintf("2 Slices of Pepperoni Pizza, French Fries, Cooler Drink"),
		Status:       ProductStatusOnShelf,
	}
	peppizza.DisplayName = properFoodName("pepperoni pizza meal")
	peppizza.DescriptionHTML = strings.Replace(peppizza.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, peppizza)

	// cheesepizza
	cheesepizza := Product{
		SKU:          "P1",
		Name:         fmt.Sprintf("cheese-pizza"),
		PriceInCents: 730,
		Description:  fmt.Sprintf("2 Slices of Cheese Pizza, French Fries, Cooler Drink"),
		Status:       ProductStatusOnShelf,
	}
	cheesepizza.DisplayName = properFoodName("cheese pizza meal")
	cheesepizza.DescriptionHTML = strings.Replace(cheesepizza.Description, ", ", "<br />", -1)
	allFoods = append(allFoods, cheesepizza)

	// coke
	coke := Product{
		SKU:             "D0",
		Name:            "coca-cola",
		DisplayName:     "Coca-Cola",
		PriceInCents:    0,
		Description:     "Coca-Cola",
		DescriptionHTML: "Coca-Cola",
		Status:          ProductStatusAddon,
	}
	allFoods = append(allFoods, coke)

	// sprite
	sprite := Product{
		SKU:             "D1",
		Name:            "sprite",
		DisplayName:     "Sprite",
		PriceInCents:    0,
		Description:     "Sprite",
		DescriptionHTML: "Sprite",
		Status:          ProductStatusAddon,
	}
	allFoods = append(allFoods, sprite)

	// dietcoke
	dietcoke := Product{
		SKU:             "D2",
		Name:            "dietcoke",
		DisplayName:     "Diet Coke",
		PriceInCents:    0,
		Description:     "Diet Coke",
		DescriptionHTML: "Diet Coke",
		Status:          ProductStatusAddon,
	}
	allFoods = append(allFoods, dietcoke)

	// orange
	orange := Product{
		SKU:             "D3",
		Name:            "orange",
		DisplayName:     "Orange Juice",
		PriceInCents:    0,
		Description:     "Orange Juice",
		DescriptionHTML: "Orange Juice",
		Status:          ProductStatusAddon,
	}
	allFoods = append(allFoods, orange)

	// mrpibbs
	mrpibbs := Product{
		SKU:             "D4",
		Name:            "mrpibbs",
		DisplayName:     "Mr Pibbs",
		PriceInCents:    0,
		Description:     "Mr Pibbs",
		DescriptionHTML: "Mr Pibbs",
		Status:          ProductStatusAddon,
	}
	allFoods = append(allFoods, mrpibbs)

	var err error
	for i := range allFoods {
		err = allFoods[i].Create()
		if err != nil {
			logger.Fatal("error while generating sample data:", err)
		}
	}

}
