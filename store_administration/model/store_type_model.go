package model

import "github.com/vukovlevi/netstore/store_administration/db"

type StoreType struct {
    Id int `json:"id"`
    Name string `json:"name"`
}

func GetAllStoreType() ([]StoreType, error) {
    rows, err := db.DB.Query("SELECT id, name FROM store_type")
    if err != nil {
        return []StoreType{}, err
    }

    storeTypes := make([]StoreType, 0)
    for rows.Next() {
        storeType := StoreType{}
        err = rows.Scan(&storeType.Id, &storeType.Name)
        if err != nil {
            return []StoreType{}, err
        }
        storeTypes = append(storeTypes, storeType)
    }
    return storeTypes, nil
}
