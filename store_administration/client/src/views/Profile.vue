<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import Feedback from "../components/Feedback.vue";
import type { Feedback as TFeedback } from "../types/Feedback";
import type { User } from "../types/User";
import { router } from "../router";

const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);
const user: Ref<User | null, User | null> = ref(null);

async function getUser() {
  try {
    const resp = await fetch("/api/user");
    const data = await resp.json();

    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    user.value = data as User;
    feedback.value = null;
  } catch (err) {
    console.error(err);
    feedback.value = {
      type: "error",
      message:
        "Ismeretlen hiba miatt nem sikerült betölteni a profilod adatait!",
    };
  }
}

async function logout() {
  try {
    const resp = await fetch("/api/logout")
    const data = await resp.json()
    if (data.error) {
      feedback.value = { type: "error", message: data.error }
      return
    }

    feedback.value = null
    router.push("/login")
  } catch (err) {
    console.error(err)
    feedback.value = { type: "error", message: "Ismeretlen okok miatt nem sikerült kijelentkeztetni!" }
  }
}

onMounted(() => {
  getUser();
});
</script>

<template>
  <div class="container mx-auto max-w-2xl px-4 py-10 sm:px-6 lg:px-8 mt-[5rem]">
    <div class="space-y-8">
      <div>
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white">Profil</h2>
      </div>

      <div
        class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-700 dark:bg-background-dark/50">
        <Feedback v-if="feedback != null" :feedback="feedback" />
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
              Vezetéknév
            </p>
            <p class="mt-1 text-base font-semibold text-gray-800 dark:text-gray-100">
              {{ user?.lastname || "-" }}
            </p>
          </div>

          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
              Keresztnév
            </p>
            <p class="mt-1 text-base font-semibold text-gray-800 dark:text-gray-100">
              {{ user?.firstname || "-" }}
            </p>
          </div>

          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
              Felhasználónév
            </p>
            <p class="mt-1 text-base font-semibold text-gray-800 dark:text-gray-100">
              {{ user?.username || "-" }}
            </p>
          </div>

          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
              Email
            </p>
            <p class="mt-1 text-base font-semibold text-gray-800 dark:text-gray-100">
              {{ user?.email || "-" }}
            </p>
          </div>

          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
              Telefonszám
            </p>
            <p class="mt-1 text-base font-semibold text-gray-800 dark:text-gray-100">
              {{ user?.phoneNumber || "-" }}
            </p>
          </div>

          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
              Rang
            </p>
            <p class="mt-1 text-base font-semibold text-gray-800 dark:text-gray-100">
              {{ user?.role || "-" }}
            </p>
          </div>

          <div>
            <RouterLink to="/password-change"
              class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
              type="submit">
              Jelszó megváltoztatása
            </RouterLink>
          </div>
          <div>
            <button type="button" @click="logout"
              class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark">
              Kijelentkezés
            </button>
          </div>
        </div>
      </div>

      <div class="flex justify-end">
        <RouterLink to="/"
          class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark">
          Vissza
        </RouterLink>
      </div>
    </div>
  </div>
</template>
