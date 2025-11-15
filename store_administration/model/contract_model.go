package model

import (
	"database/sql"
	"time"
)

type Contract struct {
	Id             int    `json:"id"`
	UserName       string `json:"userName"`
	UserId         int    `json:"userId"`
	ContractType   string `json:"contractType"`
	ContractTypeId int    `json:"contractTypeId"`
	Salary         int    `json:"salary"`
	Filename       string `json:"filename"`
	StartsAt       time.Time `json:"startsAt"`
	EndsAt 		   sql.NullTime `json:"endsAt"`
	DeletedAt 	   sql.NullTime `json:"deleted_at"`
}
