package gorm

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// GlobalVar holds the global variables. There can only be one
type GlobalVar struct {
	gorm.Model
	CurrentOrderTagNumber int
}

// FirstOrCreate : first or create
func (g *GlobalVar) FirstOrCreate() error {
	_, err := GetGlobalVar()
	if err == gorm.ErrRecordNotFound {
		return DB.Create(g).Error
	}
	return nil
}

// Save :
func (g *GlobalVar) Save() error {
	if g.ID == 0 {
		return errors.New("id is zero")
	}
	return DB.Save(g).Error
}

// GetGlobalVar : get the global var
func GetGlobalVar() (*GlobalVar, error) {
	var g GlobalVar
	err := DB.Where("id = 1").Last(&g).Error
	return &g, err
}

func initGlobalVar() error {
	var g GlobalVar
	g.CurrentOrderTagNumber = 0
	return g.FirstOrCreate()
}

// GetNextOrderTag : to generate an order tag when it's placed
func (g *GlobalVar) GetNextOrderTag() (int, error) {
	g.CurrentOrderTagNumber++
	err := g.Save()
	if err != nil {
		return 0, err
	}
	return g.CurrentOrderTagNumber, err
}