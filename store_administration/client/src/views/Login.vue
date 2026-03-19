<script setup lang="ts">
import { ref } from "vue";
import { router } from "../router/index.ts";

const username = ref("");
const password = ref("");
const showPassword = ref(false);
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
    setTimeout(() => router.push("/"), 100);
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

      <form class="space-y-6" @submit.prevent="login">
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
          <div class="relative">
            <input
              class="mt-1 block w-full px-4 py-3 pr-12 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
              id="password"
              name="password"
              placeholder="Adja meg jelszavát"
              required
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
            />

            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute inset-y-0 right-3 flex items-center text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
            >
              <!-- Eye (hidden) -->
              <svg
                v-if="!showPassword"
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                />
              </svg>

              <!-- Eye (visible) -->
              <svg
                v-else
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7a9.956 9.956 0 012.042-3.368M6.223 6.223A9.956 9.956 0 0112 5c4.478 0 8.268 2.943 9.542 7a9.97 9.97 0 01-4.043 5.132M15 12a3 3 0 00-3-3"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M3 3l18 18"
                />
              </svg>
            </button>
          </div>
        </div>

        <div>
          <button
            class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-bold text-white bg-primary hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary focus:ring-offset-background-light dark:focus:ring-offset-background-dark"
            type="submit"
          >
            Bejelentkezés
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
