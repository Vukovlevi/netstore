package model

import (
	"database/sql"
	"errors"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type OpenHour struct {
	Id      int `json:"id"`
	OpensAt string `json:"opensAt"`
	ClosesAt string `json:"closesAt"`
	SpecialDate sql.NullTime `json:"specialDate,omitempty,omitzero"`
	WeekDayIds []int `json:"weekDayIds,omitempty,omitzero"`
	WeekDays []string `json:"weekDays,omitempty,omitzero"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}

func GetOpenHours() ([]OpenHour, error) {
	rows, err := db.DB.Query("SELECT id, opens_at, closes_at, special_date, deleted_at FROM open_hour WHERE deleted_at IS NULL")
	if err != nil {
		return []OpenHour{}, err
	}

	openHours := []OpenHour{}
	for rows.Next() {
		openHour := OpenHour{}
		err = rows.Scan(&openHour.Id, &openHour.OpensAt, &openHour.ClosesAt, &openHour.SpecialDate, &openHour.DeletedAt)
		if err != nil {
			return openHours, err
		}

		weekDays, err := getWeekDaysForOpenHour(openHour.Id)
		if err != nil {
			return openHours, err
		}
		openHour.WeekDays = weekDays
		openHours = append(openHours, openHour)
	}

	return openHours, nil
}

func getWeekDaysForOpenHour(openHourId int) ([]string, error) {
	rows, err := db.DB.Query("SELECT name FROM open_day INNER JOIN week_day ON open_day.week_day_id = week_day.id WHERE open_hour_id = ? AND deleted_at IS NULL", openHourId)
	if err != nil {
		return []string{}, err
	}

	weekDays := []string{}
	for rows.Next() {
		weekDay := ""
		err = rows.Scan(&weekDay)
		if err != nil {
			return weekDays, err
		}
		weekDays = append(weekDays, weekDay)
	}

	return weekDays, nil
}

func (o *OpenHour) InsertNewOpenHour() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO open_hour (opens_at, closes_at, special_date) VALUES (?, ?, ?)", o.OpensAt, o.ClosesAt, o.SpecialDate)
	if err != nil {
		return err
	}

	id := 0
	row := tx.QueryRow("SELECT id FROM open_hour ORDER BY id DESC LIMIT 1")
	err = row.Scan(&id)
	if err != nil {
		return err
	}
	o.Id = id

	err = insertWeekDaysForOpenHour(o.Id, o.WeekDayIds, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (o *OpenHour) UpdateOpenHour() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE open_hour SET opens_at = ?, closes_at = ?, special_date = ? WHERE id = ?", o.OpensAt, o.ClosesAt, o.SpecialDate, o.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM open_day WHERE open_hour_id = ?", o.Id)
	if err != nil {
		return err
	}

	err = insertWeekDaysForOpenHour(o.Id, o.WeekDayIds, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertWeekDaysForOpenHour(openHourId int, weekDayIds []int, tx *sql.Tx) error {
	for _, weekDayId := range weekDayIds {
		_, err := tx.Exec("INSERT INTO open_day (open_hour_id, week_day_id) VALUES (?, ?)", openHourId, weekDayId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OpenHour) DeleteOpenHour() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE open_hour SET deleted_at = NOW() WHERE id = ?", o.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE open_day SET deleted_at = NOW() WHERE open_hour_id = ?", o.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (o *OpenHour) ValidateInsert() error {
	if o.OpensAt == "" || o.ClosesAt == "" {
		return errors.New("A nyitvatartási idő adatai nem megfelelők!")
	}
	return nil
}

func (o *OpenHour) ValidateUpdate() error {
	if o.Id == 0 || o.OpensAt == "" || o.ClosesAt == "" {
		return errors.New("A nyitvatartási idő frissítéséhez az adatok nem megfelelők!")
	}
	return nil
}

func (o *OpenHour) ValidateDelete() error {
	if o.Id == 0 {
		return errors.New("A nyitvatartási idő törléséhez az azonosító megadása kötelező!")
	}
	return nil
}