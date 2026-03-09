<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import type { Feedback as TFeedback } from "../types/Feedback";
import { StoreDetailClass, type StoreDetail } from "../types/StoreDetail";
import type { StoreType } from "../types/StoreType";
import Feedback from "../components/Feedback.vue";
import Modal from "../components/Modal.vue";
import { router } from "../router";

const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);
const connectionState: Ref<TFeedback | null, TFeedback | null> = ref(null);

const storeDetail = ref(new StoreDetailClass());
const psk = ref("");
let oldStoreDetail: StoreDetail;
const isModalOpen = ref(false);
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
    oldStoreDetail = data as StoreDetail;
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
    oldStoreDetail = storeDetail.value.toStoreDetail();
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
  getConnectionState();
});

function cancel() {
  isModalOpen.value = false;
}

function confirm() {
  isModalOpen.value = false;
  router.push("/");
}

async function connect() {
  try {
    const resp = await fetch("/api/connect", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        ipAddress: storeDetail.value.centralServerAddress,
        port: storeDetail.value.centralServerPort.toString(),
        psk: psk.value,
      }),
    });
    const data = await resp.json();

    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      connectionState.value = null;
      return;
    }

    feedback.value = { type: "success", message: data.message };
  } catch (err) {
    console.error(err);
    feedback.value = {
      type: "error",
      message:
        "Ismeretlen hiba miatt nem sikerült csatlakozni a központi szerverhez!",
    };
  }

  connectionState.value = null;
}

async function getConnectionState() {
  try {
    const resp = await fetch("/api/connect");
    const data = await resp.json();
    connectionState.value = { type: "info", message: data.message };
  } catch (err) {
    console.error(err);
    feedback.value = {
      type: "error",
      message:
        "Ismeretlen hiba miatt nem sikerült lekérdezni, hogy a rendszer csatlakozik-e a központi szerverhez.",
    };
  }
}
</script>

<template>
  <Modal
    v-if="isModalOpen"
    title="Biztosan vissza akar lépni?"
    message="Az üzlet adatainak nem mentett módosításai elvesznek!"
    confirm-text="Igen, visszalépek"
    @cancel="cancel"
    @confirm="confirm"
  />

  <div
    class="min-h-screen bg-gray-50 dark:bg-background-dark flex items-start justify-center px-4 py-6 sm:py-10 overflow-y-auto mt-[3rem] lg:overflow-y-hidden"
  >
    <div class="w-full max-w-3xl">
      <Feedback v-if="connectionState != null" :feedback="connectionState" />

      <form
        class="bg-white dark:bg-background-dark/70 rounded-xl shadow-lg overflow-hidden"
        @submit.prevent="saveStoreDetail"
      >
        <!-- Header -->
        <div class="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
            Üzlet adatai
          </h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Adja meg az üzlet és a központi szerver adatait
          </p>
        </div>

        <!-- Content -->
        <div class="p-6 space-y-6">
          <Feedback v-if="feedback != null" :feedback="feedback" />

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Address -->
            <div class="md:col-span-2">
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Üzlet címe*
              </label>
              <input
                type="text"
                v-model="storeDetail.address"
                placeholder="Adja meg az üzlet címét"
                class="mt-1 w-full rounded-md border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400"
                required
              />
            </div>

            <!-- Store type -->
            <div>
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Üzlet típusa*
              </label>
              <select
                v-model="storeDetail.storeTypeId"
                class="mt-1 w-full rounded-md border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white"
              >
                <option value="0">Válasszon üzlet típust</option>
                <option v-for="storeType in storeTypes" :value="storeType.id">
                  {{ storeType.name }}
                </option>
              </select>
            </div>

            <!-- Port -->
            <div>
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Központi szerver portja*
              </label>
              <input
                type="number"
                v-model="storeDetail.centralServerPort"
                placeholder="Port"
                class="mt-1 w-full rounded-md border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white"
                required
              />
            </div>

            <!-- Central server -->
            <div class="md:col-span-2">
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Központi szerver címe*
              </label>
              <input
                type="text"
                v-model="storeDetail.centralServerAddress"
                placeholder="https://example.com"
                class="mt-1 w-full rounded-md border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400"
                required
              />
            </div>

            <!-- PSK -->
            <div class="md:col-span-2">
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Előre megbeszélt kulcs (PSK)
              </label>
              <input
                type="text"
                v-model="psk"
                placeholder="Kulcs"
                class="mt-1 w-full rounded-md border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400"
                required
              />
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div
          class="bg-gray-50 dark:bg-background-dark/40 px-6 py-4 flex flex-col sm:flex-row gap-3 sm:justify-end"
        >
          <button
            type="button"
            class="rounded-md bg-green-600 px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-green-700 dark:bg-green-500 dark:hover:bg-green-600"
            @click="connect"
          >
            Csatlakozás
          </button>

          <button
            type="button"
            class="rounded-md bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-gray-300 hover:bg-gray-100 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600"
            @click="
              () => {
                if (!storeDetail.compare(oldStoreDetail)) isModalOpen = true;
                else confirm();
              }
            "
          >
            Mégse
          </button>

          <button
            type="submit"
            class="rounded-md bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90"
          >
            Mentés
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
