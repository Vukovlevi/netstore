<script setup lang="ts">
import type { Ref } from "vue";
import { onMounted, ref } from "vue";

import UserTable from "../components/users/UserTable.vue";
import SearchBar from "../components/SearchBar.vue";
import UserData from "../components/users/UserData.vue";
import type { User } from "../types/User.ts";
import type { Role } from "../types/Role.ts";
import Feedback from "../components/Feedback.vue";
import type { FeedbackType, Feedback as TFeedback } from "../types/Feedback.ts";

let users: User[] = [];
const filteredUsers: Ref<User[], User[]> = ref([]);
let roles: Role[] = [];

const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null)

const currentUser: Ref<User | null, User | null> = ref(null);
const mode: Ref<"all" | "single", "all" | "single"> = ref("all");

async function getUsers() {
  try {
    const resp = await fetch("/api/all-user");
    const data = await resp.json();

    if (data.error) {
      feedback.value = {type: "error", message: (data.error as string)};
      return;
    }

    feedback.value = null;
    users = data as User[];
    filteredUsers.value = users;
  } catch (err) {
    feedback.value = {type: "error", message: "Ismeretlen hiba miatt nem sikerült lekérni a felhasználókat!"};
    console.error(err);
  }
}

async function getRoles() {
  try {
    const resp = await fetch("/api/role");
    const data = await resp.json();

    if (data.error) {
      feedback.value = {type: "error", message: (data.error as string)};
      return;
    }

    feedback.value = null
    roles = data as Role[];
  } catch (err) {
    feedback.value = {type: "error", message: "Ismeretlen hiba miatt nem sikerült lekérni a rangokat (csak felvitelnél és módosításnál jelent problémát)!"};
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
  feedback.value = null;
  mode.value = "single"
}

function handleFeedback(type: FeedbackType, msg: string, user: User | null, isUpdate: boolean) {
  feedback.value = {type: type, message: msg}
  console.log(feedback.value)
  if (user == null) return;
  
  if (isUpdate) updateUser(user);
  else createNewUser(user);
}

function createNewUser(user: User) {
  users.push(user);
  filteredUsers.value = users;
}

function updateUser(user: User) {
  const idx = users.findIndex(x => x.id == user.id);
  users[idx] = user;
}

async function deleteUser(userId: Number) {
  try {
    const resp = await fetch("/api/user", {
      method: "DELETE",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({id: userId}),
    })

    if (!resp.ok) {
      feedback.value = {type: "error", message: "Felhasználó törlése sikertelen!"};
      return;
    }
    
    users = users.filter(x => x.id != userId);
    filteredUsers.value = users;
    feedback.value = {type: "success", message: "Felhasználó törlése sikeres!"};
  } catch (err) {
    feedback.value = {type: "error", message: "Ismeretlen hiba miatt nem sikerült törölni a felhasználót!"};
    console.error(err);
  }
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
        @click="() => {mode = 'single'; feedback = null;}"
      >
        Felhasználó felvitele
      </button>
    </div>
    <SearchBar
      search-item="Felhasználók"
      @search="search"
      v-if="mode == 'all'"
    />
    <Feedback v-if="feedback != null" :feedback="feedback"/>

    <div
      class="overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800"
    >
      <UserTable :users="filteredUsers" v-if="mode == 'all'" @modify="(user: User) => modifyUser(user)" @delete="deleteUser" />
      <UserData
        v-else
        :user="currentUser"
        :roles="roles"
        @feedback="handleFeedback"
        @back="() => {mode = 'all'; feedback = null; currentUser = null;}"
      />
    </div>
  </div>
</template>
