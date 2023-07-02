package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Employee struct {
	ID        uint32    `gorm:"primary_key,auto_increment" json:"id"`
	Name      string    `gorm:"size:255,not null" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"CreatedAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt"`
}

func (e *Employee) Prepare() {
	e.ID = 0
	e.Name = html.EscapeString(strings.TrimSpace(e.Name))
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

func (e *Employee) Validate() error {
	if e.Name == "" {
		return errors.New("Name is invalid")
	}
	return nil
}

func (e *Employee) FindAllEmployees(db *gorm.DB) (*[]Employee, error) {
	var err error
	employees := []Employee{}
	err = db.Debug().Model(&Employee{}).Limit(100).Find(&employees).Error
	if err != nil {
		return &[]Employee{}, err
	}

	return &employees, err
}

func (e *Employee) FindEmployeeByID(db *gorm.DB, uid uint32) (*Employee, error) {
	// var err error
	err := db.Debug().Model(&Employee{}).Where("id=?", uid).Take(&e).Error

	if err != nil {
		return &Employee{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Employee{}, errors.New("user not found")
	}
	return e, err
}

func (e *Employee) CreateEmployee(db *gorm.DB) (*Employee, error) {
	err := db.Debug().Create(&e).Error

	if err != nil {
		return &Employee{}, err
	}
	return e, nil
}

func (e *Employee) UpdateEmployee(db *gorm.DB, uid int32) (*Employee, error) {
	fmt.Println(e)
	// verification
	err := db.Debug().Model(&Employee{}).Where("id=?", uid).Take(&Employee{}).UpdateColumns(
		map[string]interface{}{
			"name":      e.Name,
			"updatedAt": time.Now(),
		},
	)

	if err.Error != nil {
		return &Employee{}, err.Error
	}
	err1 := db.Debug().Model(&Employee{}).Where("id=?", uid).Take(&Employee{}).Error

	if err1 != nil {
		return &Employee{}, err1
	}

	return e, nil
}

func (e *Employee) DeleteEmployee(db *gorm.DB, uid int32) (int64, error) {
	db = db.Debug().Model(&Employee{}).Where("id=?", uid).Take(&Employee{}).Delete(&Employee{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
