<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import type { Feedback as TFeedback } from "../types/Feedback";
import { StoreDetailClass, type StoreDetail } from "../types/StoreDetail";
import type { StoreType } from "../types/StoreType";
import Feedback from "../components/Feedback.vue";

const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);

const storeDetail = ref(new StoreDetailClass());
const storeTypes: Ref<StoreType[], StoreType[]> = ref([]);

async function getStoreDetail() {
  try {
    const resp = await fetch("/api/store-detail");
    const data = await resp.json();

    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    storeDetail.value = new StoreDetailClass(data as StoreDetail);
  } catch (err) {
    feedback.value = {
      type: "error",
      message: "Ismeretlen hiba történt az üzlet adatainak lekérésekor!",
    };
    console.error(err);
  }
}

async function getStoreTypes() {
  try {
    const resp = await fetch("/api/store-type");
    const data = await resp.json();

    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    storeTypes.value = data as StoreType[];
  } catch (err) {
    feedback.value = {
      type: "error",
      message: "Ismeretlen hiba történt az üzlettípusok lekérésekor!",
    };
    console.error(err);
  }
}

async function saveStoreDetail() {
  const { message, valid } = validate();
  if (!valid) {
    feedback.value = { type: "warning", message: message };
    return;
  }

  try {
    const resp = await fetch("/api/store-detail", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(storeDetail.value.toStoreDetail()),
    });

    const data = await resp.json();
    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    feedback.value = {
      type: "success",
      message: "Az üzlet adatainak mentése sikeres!",
    };
  } catch (err) {
    feedback.value = {
      type: "error",
      message: "Az üzlet adatainak mentése közben ismeretlen hiba lépett fel!",
    };
    console.error(err);
  }
}

function validate(): { message: string; valid: boolean } {
  if (
    storeDetail.value.address == "" ||
    storeDetail.value.centralServerAddress == "" ||
    storeDetail.value.storeTypeId == 0
  ) {
    return {
      message: "A *-gal jelölt mezők kitöltése kötelező!",
      valid: false,
    };
  }

  if (
    storeDetail.value.centralServerPort < 1024 ||
    storeDetail.value.centralServerPort > 65535 ||
    storeDetail.value.centralServerPort !=
      parseInt(storeDetail.value.centralServerPort.toString())
  ) {
    return {
      message:
        "A központi szerver portja nem valid (1025-65535 közötti egész szám kell legyen!)",
      valid: false,
    };
  }

  return { message: "", valid: true };
}

onMounted(() => {
  getStoreDetail();
  getStoreTypes();
});
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-background-dark"
  >
    <form
      class="w-full max-w-md space-y-6 p-6 rounded-lg shadow-md bg-white dark:bg-background-dark/70"
      @submit.prevent="saveStoreDetail"
    >
      <div class="space-y-8">
        <div>
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
            Üzlet adatai
          </h2>
        </div>
      </div>

      <Feedback v-if="feedback != null" :feedback="feedback" />

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="current_password"
        >
          Üzlet címe*
        </label>
        <div class="mt-1">
          <input
            id="current_password"
            type="text"
            placeholder="Adja meg az üzlet címét"
            v-model="storeDetail.address"
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
          Központi szerver címe*
        </label>
        <div class="mt-1">
          <input
            id="new_password"
            type="text"
            placeholder="Adja meg a központi szerver címét"
            v-model="storeDetail.centralServerAddress"
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
          Központi szerver portja*
        </label>
        <div class="mt-1">
          <input
            id="confirm_password"
            type="number"
            placeholder="Adja meg a központi szerver portját"
            v-model="storeDetail.centralServerPort"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            required
          />
        </div>
      </div>

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="store_type"
        >
          Üzlet típusa*
        </label>
        <div class="mt-1">
          <select
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:focus:border-primary dark:focus:ring-primary"
            id="store_type"
            v-model="storeDetail.storeTypeId"
          >
            <option value="0">Válasszon üzlet típust</option>
            <option v-for="storeType in storeTypes" :value="storeType.id">
              {{ storeType.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="flex justify-end gap-3 pt-4">
        <RouterLink
          to="/"
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
