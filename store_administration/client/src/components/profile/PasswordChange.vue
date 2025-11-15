<script setup lang="ts">
import { reactive, type Ref, ref } from "vue";
import Feedback from "../Feedback.vue";
import type { Feedback as TFeedback } from "../../types/Feedback";
import { router } from "../../router";

const password = reactive({
  oldPassword: "",
  newPassword: "",
  newConfirm: "",
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
        <div class="mt-1">
          <input
            id="current_password"
            type="password"
            placeholder="Adja meg jelenlegi jelszavát"
            v-model="password.oldPassword"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            required
          />
        </div>
      </div>

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="new_password"
        >
          Új jelszó*
        </label>
        <div class="mt-1">
          <input
            id="new_password"
            type="password"
            placeholder="Adja meg az új jelszót"
            v-model="password.newPassword"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            required
          />
        </div>
      </div>

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="confirm_password"
        >
          Új jelszó megerősítése*
        </label>
        <div class="mt-1">
          <input
            id="confirm_password"
            type="password"
            placeholder="Erősítse meg az új jelszót"
            v-model="password.newConfirm"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            required
          />
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
