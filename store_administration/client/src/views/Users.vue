<script setup lang="ts">
import type { Ref } from "vue";
import { onMounted, ref } from "vue";

import UserTable from "../components/users/UserTable.vue";
import SearchBar from "../components/SearchBar.vue";
import type { User } from "../types/User";

let users: User[] = [];
const filteredUsers: Ref<User[], User[]> = ref([]);
const isError = ref(false);
const errorMessage = ref("");

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

onMounted(() => {
  getUsers();
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
      >
        Felhasználó felvitele
      </button>
    </div>
    <SearchBar search-item="Felhasználók" @search="search" />
    <div
      v-if="isError"
      class="p-4 text-sm rounded-lg border border-red-400 bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800 mb-3"
      role="alert"
    >
      {{ errorMessage }}
    </div>
    <div
      class="overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800"
    >
      <UserTable :users="filteredUsers" />
    </div>
  </div>
</template>
