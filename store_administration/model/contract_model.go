package model

import (
	"database/sql"
	"time"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type ContractDay struct {
	Id int `json:"id"`
	ContractId int `json:"contractId,omitempty,omitzero"`
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
	Filename       string `json:"filename"`
	StartsAt       time.Time `json:"startsAt"`
	EndsAt 		   sql.NullTime `json:"endsAt"`
	DeletedAt 	   sql.NullTime `json:"deleted_at"`
	ContractDays   []ContractDay `json:"contractDays"`
}

func GetAllContract() ([]Contract, error) {
	rows, err := db.DB.Query("SELECT contract.id, CONCAT(user.firstname, ' ', 'user.lastname'), contract_type.name, salary, starts_at, ends_at, contract.deleted_at FROM contract INNER JOIN contract_type ON contract.contract_type_id = contract_type.id INNER JOIN user ON contract.user_id = user.id WHERE contract.deleted_at IS NULL")
	if err != nil {
		return []Contract{}, err
	}

	contracts := []Contract{}
	for rows.Next() {
		contract := Contract{}
		err = rows.Scan(&contract.Id, &contract.UserName, &contract.ContractType, &contract.Salary, &contract.StartsAt, &contract.EndsAt, &contract.DeletedAt)
		if err != nil {
			return contracts, err
		}

		contractDays, err := getContractDaysForContract(contract.Id)
		if err != nil {
			return contracts, err
		}
		contract.ContractDays = contractDays

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func getContractDaysForContract(contractId int) ([]ContractDay, error) {
	rows, err := db.DB.Query("SELECT contract_day.id, starting_time, ending_time, week_day.name, contract_day.deleted_at FROM contract_day INNER JOIN week_day ON contract_day.week_day_id = week_day.id WHERE contract_id = ? AND contract_day.deleted_at IS NULL", contractId)
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