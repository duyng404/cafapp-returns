package gorm

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// GlobalVar holds the global variables. There can only be one
type GlobalVar struct {
	gorm.Model
	CurrentOrderTagNumber int
	ActiveMenuID          uint
	IsCafAppRunning       bool
	AdminTestable         bool
	FrontpageAnnouncement string
}

// FirstOrCreate : first or create
func (g *GlobalVar) FirstOrCreate() error {
	_, err := GetGlobalVar()
	if err == ErrRecordNotFound {
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
	g.ActiveMenuID = 1
	g.IsCafAppRunning = false
	g.AdminTestable = false
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

// TurnCafAppOn : set running to true
func (g *GlobalVar) TurnCafAppOn() error {
	g.IsCafAppRunning = true
	return g.Save()
}

// TurnCafAppOff : set running to false
func (g *GlobalVar) TurnCafAppOff() error {
	g.IsCafAppRunning = false
	return g.Save()
}

// IsCafAppRunning : is it running at the moment ?
func IsCafAppRunning() (bool, error) {
	g, err := GetGlobalVar()
	if err != nil {
		return false, err
	}
	return g.IsCafAppRunning, nil
}

// SetFrontpageAnnouncement ..
func (g *GlobalVar) SetFrontpageAnnouncement(s string) error {
	g.FrontpageAnnouncement = s
	return g.Save()
}
