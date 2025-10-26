package model

import "github.com/vukovlevi/netstore/store_administration/db"

type Role struct {
    Id int `json:"id"`
    Name string `json:"name"`
}

func GetAllRole() ([]Role, error) {
    rows, err := db.DB.Query("SELECT id, name FROM role")
    if err != nil {
        return []Role{}, err
    }

    roles := make([]Role, 0)
    for rows.Next() {
        role := Role{}
        err = rows.Scan(&role.Id, &role.Name)
        if err != nil {
            return []Role{}, err
        }
        roles = append(roles, role)
    }

    return roles, nil
}
