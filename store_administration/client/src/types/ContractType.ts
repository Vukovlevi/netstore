import type { DeletedAt } from "./DeletedAt";

type ContractType = {
  id: number;
  name: string;
  weeklyHours: number;
  deletedAt: DeletedAt | null;
};

class ContractTypeClass {
  id = 0;
  name = "";
  weeklyHours = 0;
  deletedAt: DeletedAt | null = null;

  constructor(contractType: ContractType | null) {
    if (contractType == null) {
      return;
    }

    this.id = contractType.id;
    this.name = contractType.name;
    this.weeklyHours = contractType.weeklyHours;
    this.deletedAt = contractType.deletedAt;
  }

  toContractType(): ContractType {
    return {
      id: this.id,
      name: this.name,
      weeklyHours: this.weeklyHours,
      deletedAt: this.deletedAt,
    };
  }

  compare(contractType: ContractType): boolean {
    return (
      contractType.name == this.name &&
      contractType.weeklyHours == this.weeklyHours
    );
  }
}

export type { ContractType };
export { ContractTypeClass };
