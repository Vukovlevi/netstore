<script setup lang="ts">
import { onMounted, ref } from "vue";
import { UserClass, type User } from "../../types/User.ts";
import type { Role } from "../../types/Role.ts";
import Modal from "../Modal.vue";

const NEW_USER = "Új felhasználó felvitele";
const UPDATE_USER = "Felhasználó módosítása";

const props = defineProps<{ user: User | null; roles: Role[] }>();
const emits = defineEmits(["feedback", "back", "contract"]);
const user = new UserClass(props.user);
if (props.user != null)
  user.roleId.value = props.roles.find(
    (role) => role.name == user.role.value,
  )!.id;
let oldUser = user.toUser();
const isUpdate = ref(false);
const isModalOpen = ref(false);
const isPasswordReset = ref(false);

function validate(): { message: string; valid: boolean } {
  if (props.user == null) return validateNewUser();
  else return validateUpdateUser();
}

function validateNewUser(): { message: string; valid: boolean } {
  if (
    user.firstname.value == "" ||
    user.lastname.value == "" ||
    user.username.value == "" ||
    user.password.value == "" ||
    user.roleId.value == 0
  )
    return {
      message: "A *-gal jelölt mezők kitöltése kötelező!",
      valid: false,
    };

  if (user.password.value != user.passwordConfirm.value)
    return { message: "A megadott jelszavak nem egyeznek!", valid: false };

  return { message: "", valid: true };
}

function validateUpdateUser(): { message: string; valid: boolean } {
  if (user.id.value == 0) {
    return {
      message:
        "A felhasználó azonosítója hiányzik, próbáld meg frissíteni az oldalt!",
      valid: false,
    };
  }

  if (
    user.firstname.value == "" ||
    user.lastname.value == "" ||
    user.username.value == "" ||
    user.roleId.value == 0
  ) {
    return {
      message: "A *-gal jelölt mezők kitöltése kötelező!",
      valid: false,
    };
  }

  return { message: "", valid: true };
}

async function saveUser() {
  const { message, valid } = validate();
  if (!valid) {
    emits("feedback", "warning", message, null, false);
    return;
  }

  try {
    const resp = await fetch("/api/user", {
      method: isUpdate.value ? "PUT" : "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(user.toUser()),
    });
    const data = await resp.json();

    if (data.error) {
      emits(
        "feedback",
        "error",
        "A mentés közben hiba történt: " + data.error,
        null,
        false,
      );
      return;
    }

    user.role.value = props.roles.find((x) => x.id == user.roleId.value)!.name;
    emits("feedback", "success", data.message, user.toUser(), isUpdate.value);
    oldUser = user.toUser();
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt a következő műveletet nem sikerült végrehatjani: " +
        isUpdate.value
        ? UPDATE_USER
        : NEW_USER,
      null,
      false,
    );
  }
}

onMounted(() => {
  if (props.user != null) {
    const role = props.roles.find((x) => x.name == props.user!.role)!;
    user.roleId.value = role.id;
    isUpdate.value = true;
  }
});

function cancel() {
  isModalOpen.value = false;
}

function confirm() {
  isModalOpen.value = false;
  emits("back");
}

async function passwordReset() {
  if (!user.password.value || !user.passwordConfirm.value) {
    emits(
      "feedback",
      "warning",
      "A jelszót és a megerősítését ki kell tölteni",
      null,
      false,
    );
    return;
  }

  if (user.password.value != user.passwordConfirm.value) {
    emits(
      "feedback",
      "warning",
      "A megadott jelszavak nem egyeznek",
      null,
      false,
    );
    return;
  }

  try {
    const resp = await fetch("/api/password-reset", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        userId: user.id.value,
        password: user.password.value,
        passwordConfirm: user.passwordConfirm.value,
      }),
    });
    const data = await resp.json();

    if (data.error) {
      emits(
        "feedback",
        "error",
        "A mentés közben hiba történt: " + data.error,
        null,
        false,
      );
      return;
    }

    emits("feedback", "success", data.message, null, false);
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt a következő műveletet nem sikerült visszaállítani a jelszót: ",
      null,
      false,
    );
  }
}

function changeToContract() {
  if (user.id.value == 0) emits("feedback", "warning", "A felhasználó azonosítója nem található. Próbáld meg frissíteni az oldalt!", null, false)
  else emits("contract")
}
</script>

<template>
  <Modal v-if="isModalOpen" title="Biztosan vissza akar lépni?" message="A nem mentett módosítások elvesznek!"
    confirm-text="Igen, visszalépek" @cancel="cancel" @confirm="confirm" />
  <div class="container mx-auto max-w-2xl px-4 py-10 sm:px-6 lg:px-8">
    <div class="space-y-8">
      <div>
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
          Dolgozó adatai
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
          {{ isUpdate ? UPDATE_USER : NEW_USER }}
        </p>
      </div>

      <form class="space-y-6" @submit.prevent="saveUser">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="firstname">Vezetéknév*</label>
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="firstname" placeholder="pl.: Kiss" type="text" v-model="user.lastname.value" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="lastname">Keresztnév*</label>
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="lastname" placeholder="pl.: Miklós" type="text" v-model="user.firstname.value" />
            </div>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            for="username">Felhasználónév*</label>
          <div class="mt-1">
            <input
              class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
              id="username" placeholder="pl. mkiss" type="text" v-model="user.username.value" />
          </div>
        </div>

        <button v-if="isUpdate"
          class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
          type="button" @click="
            () => {
              isPasswordReset = !isPasswordReset;
            }
          ">
          Jelszó visszaállítás
        </button>

        <div v-if="!isUpdate || isPasswordReset" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="password">Jelszó*</label>
            <div class="relative">
              <input
                class="mt-1 block w-full px-4 py-3 pr-12 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
                id="password" name="password" placeholder="Adja meg jelszavát" required
                :type="user.showPassword.value ? 'text' : 'password'" v-model="user.password.value" />

              <button type="button" @click="user.showPassword.value = !user.showPassword.value"
                class="absolute inset-y-0 right-3 flex items-center text-gray-500 hover:text-gray-700 dark:hover:text-gray-300">
                <!-- Eye (hidden) -->
                <svg v-if="!user.showPassword.value" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                  viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>

                <!-- Eye (visible) -->
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7a9.956 9.956 0 012.042-3.368M6.223 6.223A9.956 9.956 0 0112 5c4.478 0 8.268 2.943 9.542 7a9.97 9.97 0 01-4.043 5.132M15 12a3 3 0 00-3-3" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3l18 18" />
                </svg>
              </button>
            </div>
          </div>

          <div v-if="!isUpdate || isPasswordReset">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="password_confirm">Jelszó
              megerősítése*</label>
            <div class="relative">
              <input
                class="mt-1 block w-full px-4 py-3 pr-12 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
                id="password" name="password" placeholder="Adja meg jelszavát" required
                :type="user.showConfirmPassword.value ? 'text' : 'password'" v-model="user.passwordConfirm.value" />

              <button type="button" @click="
                user.showConfirmPassword.value =
                !user.showConfirmPassword.value
                "
                class="absolute inset-y-0 right-3 flex items-center text-gray-500 hover:text-gray-700 dark:hover:text-gray-300">
                <!-- Eye (hidden) -->
                <svg v-if="!user.showConfirmPassword.value" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>

                <!-- Eye (visible) -->
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7a9.956 9.956 0 012.042-3.368M6.223 6.223A9.956 9.956 0 0112 5c4.478 0 8.268 2.943 9.542 7a9.97 9.97 0 01-4.043 5.132M15 12a3 3 0 00-3-3" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3l18 18" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="email">Email</label>
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="email" placeholder="pl. email@example.com" type="email" v-model="user.email.value" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="phone">Telefonszám</label>
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="phone" placeholder="pl. +36 30 123 4567" type="tel" v-model="user.phoneNumber.value" />
            </div>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="position">Rang</label>
          <div class="mt-1">
            <select
              class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:focus:border-primary dark:focus:ring-primary"
              id="position" v-model="user.roleId.value">
              <option value="0">Válasszon pozíciót</option>
              <option v-for="role in props.roles" :value="role.id">
                {{ role.name }}
              </option>
            </select>
          </div>
        </div>

        <div v-if="!isUpdate" class="flex items-center">
          <input id="disable_password_change" name="disable_password_change" type="checkbox"
            class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:focus:ring-primary"
            v-model="user.passwordChanged.value" />
          <label for="disable_password_change" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
            Jelszóváltoztatás kikapcsolása
          </label>
        </div>

        <div class="flex justify-end gap-3 pt-4">
          <button v-if="isPasswordReset"
            class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
            type="button" @click="passwordReset">
            Jelszó visszaállítás mentése
          </button>
          <button
            class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
            type="button" @click="changeToContract">
            Szerződés kezelése
          </button>
          <button
            class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
            type="button" @click="
              () => {
                if (!user.compare(oldUser)) isModalOpen = true;
                else emits('back');
              }
            ">
            Vissza
          </button>
          <button
            class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
            type="submit">
            Mentés
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
