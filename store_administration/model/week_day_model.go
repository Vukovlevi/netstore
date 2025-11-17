package model

import "github.com/vukovlevi/netstore/store_administration/db"

type WeekDay struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllWeekDay() ([]WeekDay, error) {
	rows, err := db.DB.Query("SELECT id, name FROM week_day ORDER BY id")
	if err != nil {
		return []WeekDay{}, err
	}

	weekDays := []WeekDay{}
	for rows.Next() {
		weekDay := WeekDay{}
		err = rows.Scan(&weekDay.Id, &weekDay.Name)
		if err != nil {
			return weekDays, err
		}
		weekDays = append(weekDays, weekDay)
	}

	return weekDays, nil
}