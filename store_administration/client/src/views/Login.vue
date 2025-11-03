<script setup lang="ts">
import { ref } from "vue";
import { router } from "../router/index.ts";

const username = ref("");
const password = ref("");
const isError = ref(false);
const errorMessage = ref("");

async function login() {
  if (username.value == "" || password.value == "") {
    errorMessage.value = "Hiányzó felhasználónév, vagy jelszó!";
    isError.value = true;
    return;
  }
  try {
    const resp = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    });
    const data = await resp.json();

    if (data.error) {
      errorMessage.value = data.error;
      isError.value = true;
      return;
    }

    isError.value = false;
    router.push("/");
  } catch (err) {
    errorMessage.value =
      "Ismeretlen hiba miatt nem sikerült bejelentkezni. Próbáld újra!";
    isError.value = true;
    console.error(err);
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div
      class="w-full max-w-md p-8 space-y-8 bg-white dark:bg-background-dark rounded-xl shadow-lg m-4"
    >
      <div class="text-center">
        <h1
          class="text-3xl font-bold text-foreground-light dark:text-foreground-dark"
        >
          Bejelentkezés
        </h1>
        <p
          class="mt-2 text-sm text-placeholder-light dark:text-placeholder-dark"
        >
          Jelentkezzen be a fiókjába a folytatáshoz.
        </p>
      </div>

      <div
        v-if="isError"
        class="p-4 text-sm rounded-lg border border-red-400 bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <form class="space-y-6">
        <div>
          <label
            class="text-sm font-medium text-foreground-light dark:text-foreground-dark"
            for="username"
          >
            Felhasználónév
          </label>
          <input
            class="mt-1 block w-full px-4 py-3 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
            id="username"
            name="username"
            placeholder="Adja meg felhasználónevét"
            required
            type="text"
            v-model="username"
          />
        </div>

        <div>
          <label
            class="text-sm font-medium text-foreground-light dark:text-foreground-dark"
            for="password"
          >
            Jelszó
          </label>
          <input
            class="mt-1 block w-full px-4 py-3 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
            id="password"
            name="password"
            placeholder="Adja meg jelszavát"
            required
            type="password"
            v-model="password"
          />
        </div>

        <div>
          <button
            class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-bold text-white bg-primary hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary focus:ring-offset-background-light dark:focus:ring-offset-background-dark"
            type="button"
            @click="login"
          >
            Bejelentkezés
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
