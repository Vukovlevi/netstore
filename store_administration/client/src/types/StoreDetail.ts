type StoreDetail = {
  address: string;
  centralServerAddress: string;
  centralServerPort: number;
  storeTypeId: number;
};

class StoreDetailClass {
  address = "";
  centralServerAddress = "";
  centralServerPort = 0;
  storeTypeId = 0;

  constructor(storeDetail: StoreDetail | null = null) {
    if (storeDetail == null) {
      return;
    }

    this.address = storeDetail.address;
    this.centralServerAddress = storeDetail.centralServerAddress;
    this.centralServerPort = storeDetail.centralServerPort;
    this.storeTypeId = storeDetail.storeTypeId;
  }

  toStoreDetail(): StoreDetail {
    return {
      address: this.address,
      centralServerAddress: this.centralServerAddress,
      centralServerPort: this.centralServerPort,
      storeTypeId: this.storeTypeId,
    };
  }
}

export type { StoreDetail };
export { StoreDetailClass };
