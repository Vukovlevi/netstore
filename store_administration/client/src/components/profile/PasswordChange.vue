<script setup lang="ts">
import { reactive, type Ref, ref } from "vue";
import Feedback from "../Feedback.vue";
import type { Feedback as TFeedback } from "../../types/Feedback";
import { router } from "../../router";

const password = reactive({
  oldPassword: "",
  newPassword: "",
  newConfirm: "",
  showOldPassword: false,
  showNewPassword: false,
  showNewConfirm: false,
});
const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);

function validate(): boolean {
  if (password.oldPassword == "" || password.newPassword == "") {
    feedback.value = {
      type: "warning",
      message: "A *-gal jelölt mezők kitöltése kötelező!",
    };
    return false;
  }

  if (password.newPassword != password.newConfirm) {
    feedback.value = { type: "warning", message: "A két jelszó nem egyezik!" };
    return false;
  }

  if (password.oldPassword == password.newPassword) {
    feedback.value = {
      type: "warning",
      message: "Az új jelszó nem egyezhet meg a régi jelszóval!",
    };
    return false;
  }

  if (password.newPassword.length < 8) {
    feedback.value = {
      type: "warning",
      message: "A jelszó legalább 8 karakter hosszú kell hogy legyen!",
    };
    return false;
  }

  return true;
}

async function changePassword() {
  if (!validate()) return;

  try {
    const resp = await fetch("/api/password-change", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(password),
    });

    if (resp.redirected) {
      return router.push("/login");
    }

    const data = await resp.json();
    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    feedback.value = { type: "success", message: data.message };
  } catch (err) {
    feedback.value = {
      type: "error",
      message: "Ismeretlen hiba miatt nem sikerült a jelszó változtatása!",
    };
    console.error(err);
  }
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-background-dark"
  >
    <form
      class="w-full max-w-md space-y-6 p-6 rounded-lg shadow-md bg-white dark:bg-background-dark/70"
      @submit.prevent="changePassword"
    >
      <div class="space-y-8">
        <div>
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
            Jelszó változtatása
          </h2>
        </div>
      </div>

      <Feedback v-if="feedback != null" :feedback="feedback" />

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="current_password"
        >
          Jelenlegi jelszó*
        </label>
        <div class="relative">
          <input
            class="mt-1 block w-full px-4 py-3 pr-12 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
            id="password"
            name="password"
            placeholder="Adja meg jelszavát"
            required
            :type="password.showOldPassword ? 'text' : 'password'"
            v-model="password.oldPassword"
          />

          <button
            type="button"
            @click="password.showOldPassword = !password.showOldPassword"
            class="absolute inset-y-0 right-3 flex items-center text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            <!-- Eye (hidden) -->
            <svg
              v-if="!password.showOldPassword"
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
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="new_password"
        >
          Új jelszó*
        </label>
        <div class="relative">
          <input
            class="mt-1 block w-full px-4 py-3 pr-12 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
            id="password"
            name="password"
            placeholder="Adja meg jelszavát"
            required
            :type="password.showNewPassword ? 'text' : 'password'"
            v-model="password.newPassword"
          />

          <button
            type="button"
            @click="password.showNewPassword = !password.showNewPassword"
            class="absolute inset-y-0 right-3 flex items-center text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            <!-- Eye (hidden) -->
            <svg
              v-if="!password.showNewPassword"
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
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="confirm_password"
        >
          Új jelszó megerősítése*
        </label>
        <div class="relative">
          <input
            class="mt-1 block w-full px-4 py-3 pr-12 bg-input-light dark:bg-input-dark border border-border-light dark:border-border-dark rounded-lg placeholder-placeholder-light dark:placeholder-placeholder-dark focus:ring-primary focus:border-primary"
            id="password"
            name="password"
            placeholder="Adja meg jelszavát"
            required
            :type="password.showNewConfirm ? 'text' : 'password'"
            v-model="password.newConfirm"
          />

          <button
            type="button"
            @click="password.showNewConfirm = !password.showNewConfirm"
            class="absolute inset-y-0 right-3 flex items-center text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            <!-- Eye (hidden) -->
            <svg
              v-if="!password.showNewConfirm"
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

      <div class="flex justify-end gap-3 pt-4">
        <RouterLink
          to="/profile"
          class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
        >
          Mégse
        </RouterLink>
        <button
          type="submit"
          class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        >
          Mentés
        </button>
      </div>
    </form>
  </div>
</template>
