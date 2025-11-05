<script setup lang="ts">
import type { Ref } from "vue";
import { onMounted, ref } from "vue";

import UserTable from "../components/users/UserTable.vue";
import SearchBar from "../components/SearchBar.vue";
import UserData from "../components/users/UserData.vue";
import type { User } from "../types/User.ts";
import type { Role } from "../types/Role.ts";

let users: User[] = [];
const filteredUsers: Ref<User[], User[]> = ref([]);
let roles: Role[] = [];

const isError = ref(false);
const errorMessage = ref("");
const isSuccess = ref(false);
const successMessage = ref("");

const currentUser: Ref<User | null, User | null> = ref(null);
const mode: Ref<"all" | "single", "all" | "single"> = ref("all");

async function getUsers() {
  try {
    const resp = await fetch("/api/all-user");
    const data = await resp.json();

    if (data.error) {
      errorMessage.value = data.error;
      isError.value = true;
      return;
    }

    isError.value = false;
    users = data as User[];
    filteredUsers.value = users;
  } catch (err) {
    errorMessage.value =
      "Ismeretlen hiba miatt nem sikerült lekérni a felhasználókat!";
    isError.value = true;
    console.error(err);
  }
}

async function getRoles() {
  try {
    const resp = await fetch("/api/role");
    const data = await resp.json();

    if (data.error) {
      errorMessage.value = data.error;
      isError.value = true;
      return;
    }

    isError.value = false;
    roles = data as Role[];
  } catch (err) {
    errorMessage.value =
      "Ismeretlen hiba miatt nem sikerült lekérni a rangokat (csak felvitelnél és módosításnál jelent problémát)!";
    isError.value = true;
    console.error(err);
  }
}

function search(searchValue: string) {
  if (searchValue == "") {
    filteredUsers.value = users;
    return;
  }

  filteredUsers.value = users.filter((x) =>
    (x.firstname + " " + x.lastname)
      .toLowerCase()
      .includes(searchValue.toLowerCase())
  );
}

function modifyUser(user: User) {
  currentUser.value = user
  mode.value = "single"
}

onMounted(() => {
  getUsers();
  getRoles();
});
</script>

<template>
  <div class="mx-auto max-w-7xl mt-[5rem]">
    <div
      class="mb-6 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center"
    >
      <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
        Felhasználók
      </h2>
      <button
        class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90"
        @click="() => (mode = 'single')"
      >
        Felhasználó felvitele
      </button>
    </div>
    <SearchBar
      search-item="Felhasználók"
      @search="search"
      v-if="mode == 'all'"
    />
    <div
      v-if="isError"
      class="p-3 text-sm rounded-lg border border-red-400 bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800 mb-3"
      role="alert"
    >
      {{ errorMessage }}
    </div>

    <div
      v-if="isSuccess"
      class="p-3 text-sm rounded-lg border border-green-400 bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-300 dark:border-green-800 mb-3"
      role="alert"
    >
      {{ successMessage }}
    </div>

    <!--
    <div
      v-if="isWarning"
      class="p-3 text-sm rounded-lg border border-amber-400 bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-800 mb-3"
      role="alert"
    >
      {{ warningMessage }}
    </div>

    <div
      v-if="isInfo"
      class="p-3 text-sm rounded-lg border border-blue-400 bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-800 mb-3"
      role="alert"
    >
      {{ infoMessage }}
    </div>
  -->

    <div
      class="overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800"
    >
      <UserTable :users="filteredUsers" v-if="mode == 'all'" @modify="(user: User) => modifyUser(user)" />
      <UserData
        :user="currentUser"
        :roles="roles"
        v-else
        @error="
          (msg) => {
            isError = true;
            errorMessage = msg;
          }
        "
        @back="() => (mode = 'all')"
        @success="
          (msg, user) => {
            isSuccess = true;
            successMessage = msg;
            users.push(user);
            filteredUsers = users;
          }
        "
      />
    </div>
  </div>
</template>
