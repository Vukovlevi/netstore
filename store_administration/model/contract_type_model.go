package model

import (
	"database/sql"
	"errors"

	"github.com/vukovlevi/netstore/store_administration/db"
)

const (
    DAYS_A_WEEK = 7
    MAX_WORK_HOURS_A_DAY = 16
    MAX_WEEKLY_HOURS = DAYS_A_WEEK * MAX_WORK_HOURS_A_DAY
)

type ContractType struct {
    Id int `json:"id"`
    Name string `json:"name"`
    WeeklyHours int `json:"weeklyHours"`
    DeletedAt sql.NullTime `json:"deletedAt"`
}

func (c *ContractType) InsertNewContractType() error {
    _, err := db.DB.Exec("INSERT INTO contract_type VALUES (NULL, ?, ?, NULL)", c.Name, c.WeeklyHours)
    return err
}

func GetAllContractType() ([]ContractType, error) {
    rows, err := db.DB.Query("SELECT id, name, weekly_hours, deleted_at FROM contract_type WHERE deleted_at IS NULL")
    if err != nil {
        return []ContractType{}, err
    }

    contractTypes := make([]ContractType, 0)
    for rows.Next() {
        contractType := ContractType{}
        err = rows.Scan(&contractType.Id, &contractType.Name, &contractType.WeeklyHours, &contractType.DeletedAt)
        if err != nil {
            return []ContractType{}, err
        }
        contractTypes = append(contractTypes, contractType)
    }
    return contractTypes, nil
}

func (c *ContractType) UpdateContractType() error {
    _, err := db.DB.Exec("UPDATE contract_type SET name = ?, weekly_hours = ? WHERE id = ?", c.Name, c.WeeklyHours, c.Id)
    return err
}

func (c *ContractType) DeleteContractType() error {
    _, err := db.DB.Exec("UPDATE contract_type SET deleted_at = NOW() WHERE id = ?", c.Id)
    return err
}

//Returns user-readable error
func (c *ContractType) ValidateInsert() error {
    if c.Name == "" || c.WeeklyHours == 0 {
        return errors.New("A szerződéstípus neve, vagy a heti munkaórák száma hiányzik!")
    }

    if c.WeeklyHours > MAX_WEEKLY_HOURS || c.WeeklyHours < 0 {
        return errors.New("A heti munkaórák száma nem a megadott intervallumban van!")
    }

    return nil
}

//Return user-readable error
func (c *ContractType) ValidateUpdate() error {
    if c.Id == 0 {
        return errors.New("A szerződéstípus frissítéséhez az azonosító megadása kötelező! Próbálja frissíteni az oldalt!")
    }
    return c.ValidateInsert()
}

//Returns user-readable error
func (c *ContractType) ValidateDelete() error {
    if c.Id == 0 {
        return errors.New("A szerződéstípus törléséhez az azonosító megadása kötelező! Próbálja frissíteni az oldalt!")
    }
    return nil
}
