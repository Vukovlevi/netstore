import type { DeletedAt, NullTime } from "./DeletedAt";

type ContractDay = {
  id: number;
  startingTime: string;
  endingTime: string;
  weekDayId: number;
  weekDay: string;
  deletedAt: DeletedAt | null;
};

type Contract = {
  id: number;
  userId: number;
  userName: string;
  contractTypeId: number;
  contractType: string;
  salary: number;
  startsAt: string;
  endsAt: NullTime;
  deletedAt: DeletedAt | null;
  contractDays: ContractDay[];
};

class ContractClass {
  id = 0;
  userId = 0;
  userName = "";
  contractTypeId = 0;
  contractType = "";
  salary = 0;
  startsAt = "";
  endsAt: NullTime = { Valid: false, Time: "0001-01-01T00:00:00Z" };
  inputEndsAt = "";
  deletedAt: DeletedAt | null = null;
  contractDays: ContractDay[] = [];

  constructor(contract: Contract | null = null) {
    if (contract == null) {
      return;
    }

    this.id = contract.id;
    if (contract.userId) this.userId = contract.userId;
    if (contract.userName) this.userName = contract.userName;
    if (contract.contractTypeId) this.contractTypeId = contract.contractTypeId;
    if (contract.contractType) this.contractType = contract.contractType;
    this.salary = contract.salary;
    this.startsAt = new Date(contract.startsAt.split("T")[0]!)
      .toISOString()
      .substring(0, 10);
    if (contract.endsAt && contract.endsAt.Valid) {
      this.endsAt = contract.endsAt;
      this.inputEndsAt = new Date(contract.endsAt.Time.split("T")[0]!)
        .toISOString()
        .substring(0, 10);
    }
    this.deletedAt = contract.deletedAt;
    this.contractDays = contract.contractDays;
  }

  toContract(): Contract {
    return {
      id: this.id,
      userId: this.userId,
      userName: this.userName,
      contractTypeId: this.contractTypeId,
      contractType: this.contractType,
      salary: this.salary,
      startsAt: new Date(this.startsAt).toISOString(),
      endsAt:
        this.inputEndsAt == ""
          ? { Valid: false, Time: "0001-01-01T00:00:00Z" }
          : { Valid: true, Time: new Date(this.inputEndsAt).toISOString() },
      deletedAt: this.deletedAt,
      contractDays: this.contractDays,
    };
  }
}

export type { Contract, ContractDay };
export { ContractClass };
