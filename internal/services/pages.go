package services

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
	"gorm.io/gorm"
)

func CreatePage(pageName string, state *common.State) error {
	page := &database.Page{
		Name:  pageName,
		Views: 0,
	}
	return state.Database.Create(page).Error
}

func PageExists(pageName string, state *common.State) bool {
	var count int64
	err := state.Database.Model(&database.Page{}).Where("name = ?", pageName).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func PageViews(pageName string, state *common.State) (int64, error) {
	var page database.Page
	query := state.Database.Model(&database.Page{}).Where("name = ?", pageName)
	result := query.First(&page)
	if result.Error != nil {
		return 0, result.Error
	}
	return page.Views, nil
}

func IncreasePageViews(pageName string, state *common.State) (int64, error) {
	if !PageExists(pageName, state) {
		return 0, CreatePage(pageName, state)
	}

	query := state.Database.Model(&database.Page{}).Where("name = ?", pageName)
	result := query.UpdateColumn("views", gorm.Expr("views + ?", 1))
	if result.Error != nil {
		return 0, result.Error
	}

	return PageViews(pageName, state)
}
