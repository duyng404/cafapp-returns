package gorm

import (
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
		Tags:         []Label{*labelMain},
	}
	allFoods = append(allFoods, chicken2)

	// chicken4
	chicken4 := Product{
		Tag:          "4C",
		Name:         "4pc-chicken",
		DisplayName:  "4 Chicken Strips",
		PriceInCents: 450,
		Tags:         []Label{*labelMain},
	}
	allFoods = append(allFoods, chicken4)

	// // burger
	// burger := Product{
	// 	Tag:          "CB",
	// 	Name:         fmt.Sprintf("burger"),
	// 	PriceInCents: 715,
	// 	Description:  fmt.Sprintf("Cheeseburger w/ lettuce & tomato, French Fries, Cooler Drink"),
	// 	Status:       ProductStatusOnShelf,
	// }
	// burger.DisplayName = properFoodName("burger meal")
	// burger.DescriptionHTML = strings.Replace(burger.Description, ", ", "<br />", -1)
	// allFoods = append(allFoods, burger)

	// // turkeysandwich
	// turkeysandwich := Product{
	// 	Tag:          "TS",
	// 	Name:         fmt.Sprintf("turkey-sandwich"),
	// 	PriceInCents: 835,
	// 	Description:  fmt.Sprintf("Turkey Sub Sandwich, French Fries, Cooler Drink **Toppings: turkey, lettuce, tomato, with mayo & mustard in packets"),
	// 	Status:       ProductStatusOnShelf,
	// }
	// turkeysandwich.DisplayName = properFoodName("turkey sandwich meal")
	// turkeysandwich.DescriptionHTML = "Turkey Sub Sandwich<br />French Fries<br />Cooler Drink<br />* Toppings: turkey, lettuce, tomato<br />* with mayo & mustard in packets"
	// allFoods = append(allFoods, turkeysandwich)

	// // veggiesandwich
	// veggiesandwich := Product{
	// 	Tag:          "VS",
	// 	Name:         fmt.Sprintf("veggie-sandwich"),
	// 	PriceInCents: 835,
	// 	Description:  fmt.Sprintf("Veggie Sub Sandwich, French Fries, Cooler Drink **Toppings include: hummus, cucumber, lettuce, tomato, green pepper, with mayo & mustard in packets"),
	// 	Status:       ProductStatusOnShelf,
	// }
	// veggiesandwich.DisplayName = properFoodName("veggie sandwich meal")
	// veggiesandwich.DescriptionHTML = "Veggie Sub Sandwich<br />French Fries<br />Cooler Drink<br />* Toppings: hummus, cucumber, lettuce, tomato, green pepper<br />* with mayo & mustard in packets"
	// allFoods = append(allFoods, veggiesandwich)

	// // peppizza
	// peppizza := Product{
	// 	Tag:          "PP",
	// 	Name:         fmt.Sprintf("pep-pizza"),
	// 	PriceInCents: 730,
	// 	Description:  fmt.Sprintf("2 Slices of Pepperoni Pizza, French Fries, Cooler Drink"),
	// 	Status:       ProductStatusOnShelf,
	// }
	// peppizza.DisplayName = properFoodName("pepperoni pizza meal")
	// peppizza.DescriptionHTML = strings.Replace(peppizza.Description, ", ", "<br />", -1)
	// allFoods = append(allFoods, peppizza)

	// // cheesepizza
	// cheesepizza := Product{
	// 	Tag:          "CP",
	// 	Name:         fmt.Sprintf("cheese-pizza"),
	// 	PriceInCents: 730,
	// 	Description:  fmt.Sprintf("2 Slices of Cheese Pizza, French Fries, Cooler Drink"),
	// 	Status:       ProductStatusOnShelf,
	// }
	// cheesepizza.DisplayName = properFoodName("cheese pizza meal")
	// cheesepizza.DescriptionHTML = strings.Replace(cheesepizza.Description, ", ", "<br />", -1)
	// allFoods = append(allFoods, cheesepizza)

	// // coke
	// coke := Product{
	// 	Tag:             "CC",
	// 	Name:            "coca-cola",
	// 	DisplayName:     "Coca-Cola",
	// 	PriceInCents:    0,
	// 	Description:     "Coca-Cola",
	// 	DescriptionHTML: "Coca-Cola",
	// 	Status:          ProductStatusAddon,
	// }
	// allFoods = append(allFoods, coke)

	// // sprite
	// sprite := Product{
	// 	Tag:             "SP",
	// 	Name:            "sprite",
	// 	DisplayName:     "Sprite",
	// 	PriceInCents:    0,
	// 	Description:     "Sprite",
	// 	DescriptionHTML: "Sprite",
	// 	Status:          ProductStatusAddon,
	// }
	// allFoods = append(allFoods, sprite)

	// // dietcoke
	// dietcoke := Product{
	// 	Tag:             "DC",
	// 	Name:            "dietcoke",
	// 	DisplayName:     "Diet Coke",
	// 	PriceInCents:    0,
	// 	Description:     "Diet Coke",
	// 	DescriptionHTML: "Diet Coke",
	// 	Status:          ProductStatusAddon,
	// }
	// allFoods = append(allFoods, dietcoke)

	// // orange
	// orange := Product{
	// 	Tag:             "OJ",
	// 	Name:            "orange",
	// 	DisplayName:     "Orange Juice",
	// 	PriceInCents:    0,
	// 	Description:     "Orange Juice",
	// 	DescriptionHTML: "Orange Juice",
	// 	Status:          ProductStatusAddon,
	// }
	// allFoods = append(allFoods, orange)

	// // mrpibbs
	// mrpibbs := Product{
	// 	Tag:             "PB",
	// 	Name:            "mrpibbs",
	// 	DisplayName:     "Mr Pibbs",
	// 	PriceInCents:    0,
	// 	Description:     "Mr Pibbs",
	// 	DescriptionHTML: "Mr Pibbs",
	// 	Status:          ProductStatusAddon,
	// }
	// allFoods = append(allFoods, mrpibbs)

	// var err error
	// for i := range allFoods {
	// 	err = allFoods[i].Create()
	// 	if err != nil {
	// 		logger.Fatal("error while generating sample data:", err)
	// 	}
	// }

}
