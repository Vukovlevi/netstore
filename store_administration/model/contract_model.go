package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type ContractDay struct {
	Id int `json:"id"`
	StartingTime string `json:"startingTime"`
	EndingTime string `json:"endingTime"`
	WeekDayId int `json:"weekDayId,omitempty,omitzero"`
	WeekDay string `json:"weekDay"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}

type Contract struct {
	Id             int    `json:"id"`
	UserName       string `json:"userName"`
	UserId         int    `json:"userId,omitempty,omitzero"`
	ContractType   string `json:"contractType"`
	ContractTypeId int    `json:"contractTypeId,omitempty,omitzero"`
	Salary         int    `json:"salary"`
	Filename       sql.NullString `json:"filename"`
	StartsAt       time.Time `json:"startsAt"`
	EndsAt 		   sql.NullTime `json:"endsAt"`
	DeletedAt 	   sql.NullTime `json:"deleted_at"`
	ContractDays   []ContractDay `json:"contractDays"`
}

func GetContractByUserId(userId int) (Contract, error) {
	row := db.DB.QueryRow("SELECT contract.id, CONCAT(user.lastname, ' ', user.firstname), contract_type.name, salary, starts_at, ends_at, file, contract.deleted_at FROM contract INNER JOIN contract_type ON contract.contract_type_id = contract_type.id INNER JOIN user ON contract.user_id = user.id WHERE user.id = ? AND contract.deleted_at IS NULL", userId)

	contract := Contract{}
	err := row.Scan(&contract.Id, &contract.UserName, &contract.ContractType, &contract.Salary, &contract.StartsAt, &contract.EndsAt, &contract.Filename, &contract.DeletedAt)
	if err != nil {
		return Contract{}, err
	}

	contractDays, err := getContractDaysForContract(contract.Id)
	if err != nil {
		return Contract{}, err
	}
	contract.ContractDays = contractDays

	return contract, nil
}

func getContractDaysForContract(contractId int) ([]ContractDay, error) {
	rows, err := db.DB.Query("SELECT contract_day.id, starting_time, ending_time, week_day.name, contract_day.deleted_at FROM contract_day INNER JOIN week_day ON contract_day.week_day_id = week_day.id WHERE contract_id = ?", contractId)
	if err != nil {
		return []ContractDay{}, err
	}

	contractDays := []ContractDay{}
	for rows.Next() {
		contractDay := ContractDay{}
		err = rows.Scan(&contractDay.Id, &contractDay.StartingTime, &contractDay.EndingTime, &contractDay.WeekDay, &contractDay.DeletedAt)
		if err != nil {
			return contractDays, err
		}
		contractDays = append(contractDays, contractDay)
	}

	return contractDays, nil
}

func (c *Contract) InsertNewContract() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

    res, err := tx.Exec("INSERT INTO contract (user_id, contract_type_id, salary, starts_at, ends_at, file) VALUES (?, ?, ?, ?, ?, ?)", c.UserId, c.ContractTypeId, c.Salary, c.StartsAt, c.EndsAt, c.Filename)
	if err != nil {
		return err
	}

    id, err := res.LastInsertId()
    if err != nil {
        return err
    }
	c.Id = int(id)

	err = insertContractDaysForContract(c.Id, c.ContractDays, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *Contract) UpdateContract() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE contract SET user_id = ?, contract_type_id = ?, salary = ?, starts_at = ?, ends_at = ?, file = ? WHERE id = ?", c.UserId, c.ContractTypeId, c.Salary, c.StartsAt, c.EndsAt, c.Filename, c.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM contract_day WHERE contract_id = ?", c.Id)
	if err != nil {
		return err
	}

	err = insertContractDaysForContract(c.Id, c.ContractDays, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertContractDaysForContract(contractId int, contractDays []ContractDay, tx *sql.Tx) error {
	for _, contractDay := range contractDays {
		_, err := tx.Exec("INSERT INTO contract_day (starting_time, ending_time, contract_id, week_day_id) VALUES (?, ?, ?, ?)", contractDay.StartingTime, contractDay.EndingTime, contractId, contractDay.WeekDayId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Contract) DeleteContract() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE contract SET deleted_at = NOW() WHERE id = ?", c.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE contract_day SET deleted_at = NOW() WHERE contract_id = ?", c.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *Contract) DeleteContractFileFromDB() error {
    _, err := db.DB.Exec("UPDATE contract SET file = NULL WHERE id = ?", c.Id)
    return err
}

func (c *Contract) ValidateInsert() error {
	if c.UserId == 0 || c.ContractTypeId == 0 || c.Salary == 0 || c.StartsAt.Equal(time.Time{}) {
		return errors.New("A szerződés feltöltéséhez hiányoznak adatok!")
	}
	return validateContractDays(c.ContractDays)
}

func validateContractDays(contractDays []ContractDay) error {
	if len(contractDays) == 0 {
		return errors.New("A szerződésnek legalább 1 munkanapot tartalmaznia kell!")
	}

	existingDays := make(map[int]bool)
	for _, contractDay := range contractDays {
		if contractDay.StartingTime == "" || contractDay.EndingTime == "" || contractDay.WeekDayId == 0 {
			return errors.New("Egy munkanap feltöltéséhez hiányoznak adatok!")
		}

		if _, ok := existingDays[contractDay.WeekDayId]; ok {
			return errors.New("Nem lehet két munkanapot a hét ugyanazon napjára rögzíteni!")
		}
		existingDays[contractDay.WeekDayId] = true
	}

	return nil
}

func (c *Contract) ValidateUpdate() error {
	if c.Id == 0 {
		return errors.New("A szerződés módosításához kötelező megadni az azonosítóját!")
	}
	return c.ValidateInsert()
}

func (c *Contract) ValidateDelete() error {
	if c.Id == 0 {
		return errors.New("A szerződés törléséhez kötelező megadni az azonosítóját!")
	}
	return nil
}
