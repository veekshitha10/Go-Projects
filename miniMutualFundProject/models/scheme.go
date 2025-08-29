package models

import (

	"gorm.io/gorm"
)

type SchemeTable struct {
	SchemeCode string `json:"scheme_code" gorm:"unique"`
	SchemeName string `json:"scheme_name"`
}

func SeedSchemeTableIfEmpty(db *gorm.DB) error {
	var count int64
	err := db.Model(&SchemeTable{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		schemes := []SchemeTable{
			{SchemeCode: "001", SchemeName: "Axis Bluechip Fund"},
			{SchemeCode: "002", SchemeName: "HDFC Mid-Cap Opportunities Fund"},
			{SchemeCode: "003", SchemeName: "SBI Small Cap Fund"},
			{SchemeCode: "004", SchemeName: "ICICI Prudential Equity & Debt Fund"},
			{SchemeCode: "005", SchemeName: "Mirae Asset Large Cap Fund"},
			{SchemeCode: "006", SchemeName: "Nippon India Growth Fund"},
			{SchemeCode: "007", SchemeName: "UTI Flexi Cap Fund"},
			{SchemeCode: "008", SchemeName: "Kotak Emerging Equity Fund"},
			{SchemeCode: "009", SchemeName: "Parag Parikh Flexi Cap Fund"},
			{SchemeCode: "010", SchemeName: "Quant Active Fund"},
			{SchemeCode: "011", SchemeName: "Canara Robeco Bluechip Equity Fund"},
			{SchemeCode: "012", SchemeName: "DSP Small Cap Fund"},
			{SchemeCode: "013", SchemeName: "Aditya Birla Sun Life Tax Relief 96"},
			{SchemeCode: "014", SchemeName: "ICICI Prudential Balanced Advantage Fund"},
			{SchemeCode: "015", SchemeName: "Tata Digital India Fund"},
			{SchemeCode: "016", SchemeName: "Axis Long Term Equity Fund"},
			{SchemeCode: "017", SchemeName: "SBI Magnum Multicap Fund"},
			{SchemeCode: "018", SchemeName: "Franklin India Prima Fund"},
			{SchemeCode: "019", SchemeName: "Motilal Oswal Nasdaq 100 Fund"},
			{SchemeCode: "020", SchemeName: "Edelweiss Balanced Advantage Fund"},
		}

		if err := db.Create(&schemes).Error; err != nil {
			return err
		}
	}

	return nil
}
