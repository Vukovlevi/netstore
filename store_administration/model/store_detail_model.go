package model

import (
	"errors"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type StoreDetail struct {
    Address string `json:"address"`
    CentralServerAddress string `json:"centralServerAddress,omitzero,omitempty"`
    CentralServerPort uint16 `json:"centralServerPort,omitzero,omitempty"`
    StoreTypeId int `json:"storeTypeId"`
    StoreTypeName string `json:"storeTypeName,omitempty,omitzero"`
}

func GetStoreDetail() (StoreDetail, error) {
    row := db.DB.QueryRow("SELECT address, central_server_address, central_server_port, store_type_id, name FROM store_detail INNER JOIN store_type ON store_detail.store_type_id = store_type.id")
    storeDetail := StoreDetail{}
    err := row.Scan(&storeDetail.Address, &storeDetail.CentralServerAddress, &storeDetail.CentralServerPort, &storeDetail.StoreTypeId, &storeDetail.StoreTypeName)
    return storeDetail, err
}

func (s *StoreDetail) UpdateStoreDetail() error {
    _, err := db.DB.Exec("UPDATE store_detail SET address = ?, central_server_address = ?, central_server_port = ?, store_type_id = ?", s.Address, s.CentralServerAddress, s.CentralServerPort, s.StoreTypeId)
    return err
}

//Returns user-readable error
func (s *StoreDetail) ValidateUpdate() error {
    if s.Address == "" || s.CentralServerAddress == "" || s.CentralServerPort == 0 || s.StoreTypeId == 0 {
        return errors.New("Az üzlet adatai hiányosak (cím, központi szerver címe, központi szerver portja, vagy az üzlet típusa)!")
    }
    return nil
}
