<script setup lang="ts">
import { onMounted, ref } from "vue";
import { UserClass, type User } from "../../types/User.ts";
import type { Role } from "../../types/Role.ts";

const NEW_USER = "Új felhasználó felvitele";
const UPDATE_USER = "Felhasználó módosítása";
let httpMethod = "POST";
const display = ref(NEW_USER);

const props = defineProps<{ user: User | null; roles: Role[] }>();
const emits = defineEmits(["feedback", "back"]);
const user = new UserClass(props.user);

//TODO: elrejteni a nem módosítható tulajdonságokat
//TODO: megfelelően kezelni a módosítást (jelenleg új felhasználóként jelenik meg a módosított is)

//TODO: átcsinálni, hogy külön kezelje a módosítást és új felvitelt
function validate(): { message: string; valid: boolean } {
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

async function saveUser() {
  const { message, valid } = validate();
  if (!valid) {
    emits("feedback", "error", message);
    return;
  }

  try {
    const resp = await fetch("/api/user", {
      method: httpMethod,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(user.toUser()),
    });
    const data = await resp.json();

    if (data.error) {
      emits("feedback", "error", "A mentés közben hiba történt: " + data.error);
      return;
    }

    user.role.value = props.roles.find((x) => x.id == user.roleId.value)!.name;
    emits("feedback", "success", data.message, user.toUser());
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt a következő műveletet nem sikerült végrehatjani: " +
        display.value
    );
  }
}

onMounted(() => {
  if (props.user != null) {
    display.value = UPDATE_USER;
    httpMethod = "PUT";
  }
});

//TODO: üzenetek törlése nézetváltáskor + feedback kezelés + vissza gomb lenyomásakor figyelmeztetés (modal) -> majd törlésnél is lehet ezt használni (hasonló mint a feedback, csak modal)
</script>

<template>
  <div class="container mx-auto max-w-2xl px-4 py-10 sm:px-6 lg:px-8">
    <div class="space-y-8">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          Dolgozó adatai
        </h1>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
          {{ display }}
        </p>
      </div>

      <form class="space-y-6" @submit.prevent="saveUser">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="firstname"
              >Vezetéknév*</label
            >
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="firstname"
                placeholder="pl.: Kiss"
                type="text"
                v-model="user.lastname.value"
              />
            </div>
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="lastname"
              >Keresztnév*</label
            >
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="lastname"
                placeholder="pl.: Miklós"
                type="text"
                v-model="user.firstname.value"
              />
            </div>
          </div>
        </div>

        <div>
          <label
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            for="username"
            >Felhasználónév*</label
          >
          <div class="mt-1">
            <input
              class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
              id="username"
              placeholder="pl. mkiss"
              type="text"
              v-model="user.username.value"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="password"
              >Jelszó*</label
            >
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="password"
                placeholder="Adja meg a jelszót"
                type="password"
                v-model="user.password.value"
              />
            </div>
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="password_confirm"
              >Jelszó megerősítése*</label
            >
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="password_confirm"
                placeholder="Erősítse meg a jelszót"
                type="password"
                v-model="user.passwordConfirm.value"
              />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="email"
              >Email</label
            >
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="email"
                placeholder="pl. email@example.com"
                type="email"
                v-model="user.email.value"
              />
            </div>
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              for="phone"
              >Telefonszám</label
            >
            <div class="mt-1">
              <input
                class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                id="phone"
                placeholder="pl. +36 30 123 4567"
                type="tel"
                v-model="user.phoneNumber.value"
              />
            </div>
          </div>
        </div>
        <div>
          <label
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            for="position"
            >Rang</label
          >
          <div class="mt-1">
            <select
              class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:focus:border-primary dark:focus:ring-primary"
              id="position"
              v-model="user.roleId.value"
            >
              <option value="0">Válasszon pozíciót</option>
              <option v-for="role in props.roles" :value="role.id">
                {{ role.name }}
              </option>
            </select>
          </div>
        </div>

        <div class="flex items-center">
          <input
            id="disable_password_change"
            name="disable_password_change"
            type="checkbox"
            class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:focus:ring-primary"
            v-model="user.passwordChanged.value"
          />
          <label
            for="disable_password_change"
            class="ml-2 block text-sm text-gray-700 dark:text-gray-300"
          >
            Jelszóváltoztatás kikapcsolása
          </label>
        </div>

        <div class="flex justify-end gap-3 pt-4">
          <button
            class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
            type="button"
            @click="() => emits('back')"
          >
            Vissza
          </button>
          <button
            class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
            type="submit"
          >
            Mentés
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
