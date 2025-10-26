package model

import (
	"database/sql"

	"github.com/vukovlevi/netstore/store_administration/db"
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
    }
    return contractTypes, nil
}

func (c *ContractType) UpdateContractType() error {
    _, err := db.DB.Exec("UPDATE contract_type SET name = ?, weekly_hours = ? WHERE id = ?", c.Name, c.WeeklyHours, c.Id)
    return err
}

func (c *ContractType) DeleteContractType() error {
    _, err := db.DB.Exec("DELETE FROM contract_type WHERE id = ?", c.Id)
    return err
}
