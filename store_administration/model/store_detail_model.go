package model

import "github.com/vukovlevi/netstore/store_administration/db"

type StoreDetail struct {
    Address string `json:"address"`
    CentralServerAddress string `json:"centralServerAddress"`
    CentralServerPort uint16 `json:"centralServerPort"`
    StoreTypeId int `json:"storeTypeId"`
}

func GetStoreDetail() (StoreDetail, error) {
    row := db.DB.QueryRow("SELECT address, central_server_address, central_server_port, store_type_id FROM store_detail")
    storeDetail := StoreDetail{}
    err := row.Scan(&storeDetail.Address, &storeDetail.CentralServerAddress, &storeDetail.CentralServerPort, &storeDetail.StoreTypeId)
    return storeDetail, err
}

func (s *StoreDetail) UpdateStoreDetail() error {
    _, err := db.DB.Exec("UPDATE store_detail SET address = ?, central_server_address = ?, central_server_port = ?, store_type_id = ?", s.Address, s.CentralServerAddress, s.CentralServerPort, s.StoreTypeId)
    return err
}
