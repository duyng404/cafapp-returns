package gorm

import (
	"cafapp-returns/logger"
	"strings"
)

func properFoodName(in string) string {
	return strings.Title(strings.Replace(in, "-", " ", -1))
}

func initData() {
	labelMain, _ := GetOrCreateLabel(ProductLabelMain)
	labelSides, _ := GetOrCreateLabel(ProductLabelSide)
	labelDrink, _ := GetOrCreateLabel(ProductLabelDrink)

	allFoods := []Product{}

	// chicken2
	chicken2 := Product{
		Tag:          "2C",
		Name:         "2pc-chicken",
		DisplayName:  "2 Chicken Strips",
		PriceInCents: 250,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, chicken2)

	// chicken4
	chicken4 := Product{
		Tag:          "4C",
		Name:         "4pc-chicken",
		DisplayName:  "4 Chicken Strips",
		PriceInCents: 450,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, chicken4)

	// burger
	burger := Product{
		Tag:          "CB",
		Name:         "burger",
		DisplayName:  "Cheeseburger",
		PriceInCents: 410,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, burger)

	// turkeysandwich
	turkeysandwich := Product{
		Tag:          "TS",
		Name:         "turkey-sandwich",
		DisplayName:  "Turkey Sandwich",
		PriceInCents: 495,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, turkeysandwich)

	// veggiesandwich
	veggiesandwich := Product{
		Tag:          "VS",
		Name:         "veggie-sandwich",
		DisplayName:  "Turkey Sandwich",
		PriceInCents: 495,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, veggiesandwich)

	// peppizza
	peppizza := Product{
		Tag:          "PP",
		Name:         "pep-pizza",
		DisplayName:  "Pepperoni Pizza",
		PriceInCents: 390,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, peppizza)

	// cheesepizza
	cheesepizza := Product{
		Tag:          "CP",
		Name:         "cheese-pizza",
		DisplayName:  "Cheese Pizza",
		PriceInCents: 390,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, cheesepizza)

	// fries
	fries := Product{
		Tag:          "FR",
		Name:         "fries",
		DisplayName:  "French Fries",
		PriceInCents: 150,
		Labels:       []Label{*labelSides},
	}
	allFoods = append(allFoods, fries)

	// chips
	chips := Product{
		Tag:          "CH",
		Name:         "chips",
		DisplayName:  "Packaged Chips",
		PriceInCents: 125,
		Labels:       []Label{*labelSides},
	}
	allFoods = append(allFoods, chips)

	// coke
	coke := Product{
		Tag:          "CC",
		Name:         "coke",
		DisplayName:  "Coca-Cola",
		PriceInCents: 165,
		Labels:       []Label{*labelDrink},
	}
	allFoods = append(allFoods, coke)

	// sprite
	sprite := Product{
		Tag:          "SP",
		Name:         "sprite",
		DisplayName:  "Sprite",
		PriceInCents: 165,
		Labels:       []Label{*labelDrink},
	}
	allFoods = append(allFoods, sprite)

	// dietcoke
	dietcoke := Product{
		Tag:          "DC",
		Name:         "diet-coke",
		DisplayName:  "Diet Coke",
		PriceInCents: 165,
		Labels:       []Label{*labelDrink},
	}
	allFoods = append(allFoods, dietcoke)

	var err error
	for i := range allFoods {
		err = allFoods[i].Create()
		if err != nil {
			logger.Fatal("error while generating initial data:", err)
		}
	}

}
