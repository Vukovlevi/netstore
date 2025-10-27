import type { DeletedAt } from "./DeletedAt";

type User = {
  id: number;
  firstname: string;
  lastname: string;
  username: string;
  phoneNumber: string;
  email: string;
  role: string;
  deletedAt: DeletedAt;
};

export type { User };
