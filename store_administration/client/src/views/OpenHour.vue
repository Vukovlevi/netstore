<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import type { FeedbackType, Feedback as TFeedback } from "../types/Feedback";
import type { OpenHour } from "../types/OpenHour";
import Feedback from "../components/Feedback.vue";
import OpenHourData from "../components/open_hours/OpenHourData.vue";
import OpenHourCard from "../components/open_hours/OpenHourCard.vue";
import Modal from "../components/Modal.vue";

const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);
const openHours: Ref<OpenHour[], OpenHour[]> = ref([]);
const currentOpenHour: Ref<OpenHour | null, OpenHour | null> = ref(null);
const mode: Ref<"all" | "single", "all" | "single"> = ref("all");
const isModalOpen = ref(false);
let toDeleteOpenHourId = 0;

async function getOpenHours() {
  try {
    const resp = await fetch("/api/open-hour");
    const data = await resp.json();

    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    feedback.value = null;
    openHours.value = data as OpenHour[];
  } catch (err) {
    console.error(err);
    feedback.value = {
      type: "error",
      message:
        "Ismeretlen hiba miatt nem sikerült lekérni a nyitvatartási időket!",
    };
  }
}

function modifyOpenHour(openHour: OpenHour) {
  currentOpenHour.value = openHour;
  feedback.value = null;
  mode.value = "single";
}

function handleFeedback(
  type: FeedbackType,
  msg: string,
  openHour: OpenHour | null,
  isUpdate: boolean
) {
  feedback.value = { type: type, message: msg };
  if (openHour == null) return;

  if (isUpdate) updateOpenHour(openHour);
  else createNewOpenHour(openHour);
}

function createNewOpenHour(openHour: OpenHour) {
  openHours.value.push(openHour);
}

function updateOpenHour(openHour: OpenHour) {
  const idx = openHours.value.findIndex((x) => x.id == openHour.id);
  openHours.value[idx] = openHour;
}

async function deleteOpenHour(openHourId: Number) {
  if (openHourId == 0) {
    feedback.value = {
      type: "warning",
      message:
        "A nyitvatartási idő azonosítója hiányzik, próbáld meg frissíteni az oldalt!",
    };
    return;
  }

  try {
    const resp = await fetch("/api/open-hour", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id: openHourId }),
    });

    if (resp.status == 204) {
      openHours.value = openHours.value.filter((x) => x.id != openHourId);
      feedback.value = {
        type: "success",
        message: "Nyitvatartási idő törlése sikeres!",
      };
      toDeleteOpenHourId = 0;
      return;
    }

    const data = await resp.json();
    if (data.error) {
      feedback.value = { type: "error", message: data.error };
      return;
    }

    feedback.value = {
      type: "error",
      message: "Nyitvatartási idő törlése sikertelen!",
    };
  } catch (err) {
    feedback.value = {
      type: "error",
      message:
        "Ismeretlen hiba miatt nem sikerült törölni a nyitvatartási időt!",
    };
    console.error(err);
  }
}

onMounted(() => {
  getOpenHours();
});

function cancel() {
  isModalOpen.value = false;
}

function confirm() {
  isModalOpen.value = false;
  deleteOpenHour(toDeleteOpenHourId);
}
</script>

<template>
  <Modal
    v-if="isModalOpen"
    title="Megerősítés"
    message="Biztosan törölni akarja a nyitvatartási időt?"
    confirm-text="Törlés"
    @cancel="cancel"
    @confirm="confirm"
  />
  <div class="mx-auto max-w-7xl mt-[5rem]">
    <div
      class="mb-6 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center"
    >
      <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
        Nyitvatartási idő
      </h2>
      <button
        v-if="mode == 'all'"
        class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90"
        @click="
          () => {
            mode = 'single';
            feedback = null;
          }
        "
      >
        Nyitvatartási idő felvitele
      </button>
    </div>
    <Feedback v-if="feedback != null" :feedback="feedback" />

    <div
      v-if="mode == 'all'"
      class="flex justify-center flex-wrap gap-6 items-center"
    >
      <OpenHourCard
        v-for="openHour in openHours"
        :openHour="openHour"
        @modify="(openHour: OpenHour) => modifyOpenHour(openHour)"
        @delete="(openHourId: number) => {toDeleteOpenHourId = openHourId; isModalOpen = true}"
      />
    </div>
    <OpenHourData
      v-else
      :openHour="currentOpenHour"
      @feedback="handleFeedback"
      @back="
        () => {
          mode = 'all';
          feedback = null;
          currentOpenHour = null;
        }
      "
    />
  </div>
</template>
