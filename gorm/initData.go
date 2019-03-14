package gorm

import (
	"cafapp-returns/logger"
)

func initData() {
	labelMain, _ := GetOrCreateLabel(ProductLabelMain)
	labelSides, _ := GetOrCreateLabel(ProductLabelSide)
	labelDrink, _ := GetOrCreateLabel(ProductLabelDrink)

	allFoods := []*Product{}

	// chicken2
	chicken2 := Product{
		Tag:          "2C",
		Name:         "2pc-chicken",
		DisplayName:  "2 Chicken Strips",
		PriceInCents: 250,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &chicken2)

	// chicken4
	chicken4 := Product{
		Tag:          "4C",
		Name:         "4pc-chicken",
		DisplayName:  "4 Chicken Strips",
		PriceInCents: 450,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &chicken4)

	// burger
	burger := Product{
		Tag:          "CB",
		Name:         "burger",
		DisplayName:  "Cheeseburger",
		PriceInCents: 410,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &burger)

	// turkeysandwich
	turkeysandwich := Product{
		Tag:          "TS",
		Name:         "turkey-sandwich",
		DisplayName:  "Turkey Sandwich",
		PriceInCents: 495,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &turkeysandwich)

	// veggiesandwich
	veggiesandwich := Product{
		Tag:          "VS",
		Name:         "veggie-sandwich",
		DisplayName:  "Turkey Sandwich",
		PriceInCents: 495,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &veggiesandwich)

	// peppizza
	peppizza := Product{
		Tag:          "PP",
		Name:         "pep-pizza",
		DisplayName:  "Pepperoni Pizza",
		PriceInCents: 390,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &peppizza)

	// cheesepizza
	cheesepizza := Product{
		Tag:          "CP",
		Name:         "cheese-pizza",
		DisplayName:  "Cheese Pizza",
		PriceInCents: 390,
		Labels:       []Label{*labelMain},
	}
	allFoods = append(allFoods, &cheesepizza)

	// fries
	fries := Product{
		Tag:          "FR",
		Name:         "fries",
		DisplayName:  "French Fries",
		PriceInCents: 150,
		Labels:       []Label{*labelSides},
	}
	allFoods = append(allFoods, &fries)

	// chips
	chips := Product{
		Tag:          "CH",
		Name:         "chips",
		DisplayName:  "Packaged Chips",
		PriceInCents: 125,
		Labels:       []Label{*labelSides},
	}
	allFoods = append(allFoods, &chips)

	// coke
	coke := Product{
		Tag:          "CC",
		Name:         "coke",
		DisplayName:  "Coca-Cola",
		PriceInCents: 165,
		Labels:       []Label{*labelDrink},
	}
	allFoods = append(allFoods, &coke)

	// sprite
	sprite := Product{
		Tag:          "SP",
		Name:         "sprite",
		DisplayName:  "Sprite",
		PriceInCents: 165,
		Labels:       []Label{*labelDrink},
	}
	allFoods = append(allFoods, &sprite)

	// dietcoke
	dietcoke := Product{
		Tag:          "DC",
		Name:         "diet-coke",
		DisplayName:  "Diet Coke",
		PriceInCents: 165,
		Labels:       []Label{*labelDrink},
	}
	allFoods = append(allFoods, &dietcoke)

	var err error
	for i := range allFoods {
		err = allFoods[i].Create()
		if err != nil {
			logger.Fatal("error while generating initial data:", err)
		}
	}

	// create menu and menu items

	// menu 1: from 7pm to 9pm
	menuOne := Menu{
		Name:        "before-9",
		DisplayName: "Before 9pm",
		Description: "",
	}
	err = menuOne.Create()
	if err != nil {
		logger.Fatal("error while generating data:", err)
	}

	// menu 2: from 9pm to 10pm
	menuTwo := Menu{
		Name:        "after-9",
		DisplayName: "After 9pm",
		Description: "After 9pm, we will be offering Packaged Chips in place of French Fries. The prices will be adjusted accordingly.",
	}
	err = menuTwo.Create()
	if err != nil {
		logger.Fatal("error while generating data:", err)
	}

	allMenuItems := []MenuItem{}

	chicken2meal1 := MenuItem{
		DisplayName:         "2pc Chicken Tender Meal",
		DisplayPriceInCents: 565,
		Description:         "2 Chicken Strips, French Fries, Cooler Drink",
		DescriptionHTML:     "2 Chicken Strips<br />French Fries<br />Cooler Drink",
		StartingMain:        chicken2,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, chicken2meal1)

	chicken4meal1 := MenuItem{
		DisplayName:         "4pc Chicken Tender Meal",
		DisplayPriceInCents: 765,
		Description:         "4 Chicken Strips, French Fries, Cooler Drink",
		DescriptionHTML:     "4 Chicken Strips<br />French Fries<br />Cooler Drink",
		StartingMain:        chicken4,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, chicken4meal1)

	burgermeal1 := MenuItem{
		DisplayName:         "Burger Meal",
		DisplayPriceInCents: 725,
		Description:         "Cheeseburger (wheat bun, lettuce, tomato), French Fries, Cooler Drink",
		DescriptionHTML:     "Cheeseburger<br />(wheat bun, lettuce, tomato)<br />French Fries<br />Cooler Drink",
		StartingMain:        burger,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, burgermeal1)

	turkeysandwichmeal1 := MenuItem{
		DisplayName:         "Turkey Sandwich Meal",
		DisplayPriceInCents: 810,
		Description:         "Turkey Sub Sandwich, French Fries, Cooler Drink. **Toppings: turkey, lettuce, tomato. **Sides: mayo, mustard.",
		DescriptionHTML:     "Turkey Sub Sandwich<br />French Fries<br />Cooler Drink<br />** Toppings: turkey, lettuce, tomato.<br />** Sides: mayo, mustard.",
		StartingMain:        turkeysandwich,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, turkeysandwichmeal1)

	veggiesandwichmeal1 := MenuItem{
		DisplayName:         "Veggie Sandwich Meal",
		DisplayPriceInCents: 810,
		Description:         "Veggie Sub Sandwich, French Fries, Cooler Drink. **Toppings: hummus, cucumber, lettuce, tomato, green bell pepper. **Sides: mayo, mustard.",
		DescriptionHTML:     "Veggie Sub Sandwich<br />French Fries<br />Cooler Drink<br />** Toppings: hummus, cucumber, lettuce, tomato, green bell pepper.<br />** Sides: mayo, mustard.",
		StartingMain:        veggiesandwich,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, veggiesandwichmeal1)

	peppizzameal1 := MenuItem{
		DisplayName:         "Pepperoni Pizza Meal",
		DisplayPriceInCents: 705,
		Description:         "2 Slices of Pepperoni Pizza, French Fries, Cooler Drink",
		DescriptionHTML:     "2 Slices of Pepperoni Pizza<br />French Fries<br />Cooler Drink",
		StartingMain:        peppizza,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, peppizzameal1)

	cheesepizzameal1 := MenuItem{
		DisplayName:         "Cheese Pizza Meal",
		DisplayPriceInCents: 705,
		Description:         "2 Slices of Cheese Pizza, French Fries, Cooler Drink",
		DescriptionHTML:     "2 Slices of Cheese Pizza<br />French Fries<br />Cooler Drink",
		StartingMain:        peppizza,
		StartingSide:        fries,
		Menu:                menuOne,
	}
	allMenuItems = append(allMenuItems, cheesepizzameal1)

	// meal2

	chicken2meal2 := MenuItem{
		DisplayName:         "2pc Chicken Tender Meal",
		DisplayPriceInCents: 565,
		Description:         "2 Chicken Strips, Packaged Chips, Cooler Drink",
		DescriptionHTML:     "2 Chicken Strips<br />Packaged Chips<br />Cooler Drink",
		StartingMain:        chicken2,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, chicken2meal2)

	chicken4meal2 := MenuItem{
		DisplayName:         "4pc Chicken Tender Meal",
		DisplayPriceInCents: 765,
		Description:         "4 Chicken Strips, Packaged Chips, Cooler Drink",
		DescriptionHTML:     "4 Chicken Strips<br />Packaged Chips<br />Cooler Drink",
		StartingMain:        chicken4,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, chicken4meal2)

	burgermeal2 := MenuItem{
		DisplayName:         "Burger Meal",
		DisplayPriceInCents: 725,
		Description:         "Cheeseburger (wheat bun, lettuce, tomato), Packaged Chips, Cooler Drink",
		DescriptionHTML:     "Cheeseburger<br />(wheat bun, lettuce, tomato)<br />Packaged Chips<br />Cooler Drink",
		StartingMain:        burger,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, burgermeal2)

	turkeysandwichmeal2 := MenuItem{
		DisplayName:         "Turkey Sandwich Meal",
		DisplayPriceInCents: 810,
		Description:         "Turkey Sub Sandwich, Packaged Chips, Cooler Drink. **Toppings: turkey, lettuce, tomato. **Sides: mayo, mustard.",
		DescriptionHTML:     "Turkey Sub Sandwich<br />Packaged Chips<br />Cooler Drink<br />** Toppings: turkey, lettuce, tomato.<br />** Sides: mayo, mustard.",
		StartingMain:        turkeysandwich,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, turkeysandwichmeal2)

	veggiesandwichmeal2 := MenuItem{
		DisplayName:         "Veggie Sandwich Meal",
		DisplayPriceInCents: 810,
		Description:         "Veggie Sub Sandwich, Packaged Chips, Cooler Drink. **Toppings: hummus, cucumber, lettuce, tomato, green bell pepper. **Sides: mayo, mustard.",
		DescriptionHTML:     "Veggie Sub Sandwich<br />Packaged Chips<br />Cooler Drink<br />** Toppings: hummus, cucumber, lettuce, tomato, green bell pepper.<br />** Sides: mayo, mustard.",
		StartingMain:        veggiesandwich,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, veggiesandwichmeal2)

	peppizzameal2 := MenuItem{
		DisplayName:         "Pepperoni Pizza Meal",
		DisplayPriceInCents: 705,
		Description:         "2 Slices of Pepperoni Pizza, Packaged Chips, Cooler Drink",
		DescriptionHTML:     "2 Slices of Pepperoni Pizza<br />Packaged Chips<br />Cooler Drink",
		StartingMain:        peppizza,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, peppizzameal2)

	cheesepizzameal2 := MenuItem{
		DisplayName:         "Cheese Pizza Meal",
		DisplayPriceInCents: 705,
		Description:         "2 Slices of Cheese Pizza, Packaged Chips, Cooler Drink",
		DescriptionHTML:     "2 Slices of Cheese Pizza<br />Packaged Chips<br />Cooler Drink",
		StartingMain:        peppizza,
		StartingSide:        chips,
		Menu:                menuTwo,
	}
	allMenuItems = append(allMenuItems, cheesepizzameal2)

	for i := range allMenuItems {
		err = allMenuItems[i].Create()
		if err != nil {
			logger.Fatal("error while generating initial data:", err)
		}
	}
}
