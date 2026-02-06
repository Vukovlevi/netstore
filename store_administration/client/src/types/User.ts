import { ref, type Ref } from "vue";
import type { DeletedAt } from "./DeletedAt";

type User = {
  id: number;
  firstname: string;
  lastname: string;
  username: string;
  password: string;
  passwordConfirm: string;
  passwordChanged: boolean;
  phoneNumber: string;
  email: string;
  role: string;
  roleId: number;
  deletedAt: DeletedAt | null;
};

class UserClass {
  id = ref(0);
  firstname = ref("");
  lastname = ref("");
  username = ref("");
  password = ref("");
  passwordConfirm = ref("");
  showPassword = ref(false);
  showConfirmPassword = ref(false);
  passwordChanged = ref(false);
  phoneNumber = ref("");
  email = ref("");
  role = ref("");
  roleId = ref(0);
  deletedAt: Ref<DeletedAt | null, DeletedAt | null> = ref(null);
  constructor(user: User | null = null) {
    if (user == null) {
      return;
    }

    this.id.value = user.id;
    this.firstname.value = user.firstname;
    this.lastname.value = user.lastname;
    this.username.value = user.username;
    this.password.value = user.password;
    this.passwordConfirm.value = user.passwordConfirm;
    this.passwordChanged.value = user.passwordChanged;
    this.phoneNumber.value = user.phoneNumber;
    this.email.value = user.email;
    this.role.value = user.role;
    this.roleId.value = user.roleId;
    this.deletedAt.value = user.deletedAt;
  }
  toUser(): User {
    return {
      id: this.id.value,
      firstname: this.firstname.value,
      lastname: this.lastname.value,
      username: this.username.value,
      password: this.password.value,
      passwordConfirm: this.passwordConfirm.value,
      passwordChanged: this.passwordChanged.value,
      phoneNumber: this.phoneNumber.value.replaceAll(" ", ""),
      email: this.email.value,
      role: this.role.value,
      roleId: this.roleId.value,
      deletedAt: this.deletedAt.value,
    };
  }

  compare(user: User): boolean {
    return (
      this.firstname.value == user.firstname &&
      this.lastname.value == user.lastname &&
      this.username.value == user.username &&
      this.phoneNumber.value == user.phoneNumber &&
      this.email.value == user.email &&
      this.roleId.value == user.roleId
    );
  }
}

export type { User };
export { UserClass };
