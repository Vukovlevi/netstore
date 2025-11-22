<script setup lang="ts">
import { ref } from "vue";
import type { OpenHour } from "../../types/OpenHour";

const CLOSED_TEXT = "Zárva";

const { openHour } = defineProps<{ openHour: OpenHour }>();
const emits = defineEmits(["modify", "delete"]);
const openTime = ref(openHour.opensAt + " - " + openHour.closesAt);
if (
  (openHour.opensAt == "00:00:00" || openHour.opensAt == "00:00") &&
  openHour.closesAt == openHour.opensAt
) {
  openTime.value = CLOSED_TEXT;
}
const openOn = ref(openHour.weekDays?.join(", "));
if (openHour.specialDate?.Valid) {
  openOn.value = new Date(openHour.specialDate.Time.split("T")[0]!)
    .toISOString()
    .substring(0, 10);
}
</script>

<template>
  <div
    class="bg-white shadow-md rounded-lg p-5 max-w-sm hover:shadow-lg transition-shadow"
  >
    <h3 class="text-lg font-semibold mb-2">{{ openTime }}</h3>
    <p class="text-gray-600 text-md mb-4">
      {{ openOn }}
    </p>

    <div class="flex gap-2">
      <button
        @click="() => emits('modify', openHour)"
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-3 rounded transition-colors"
      >
        Módosítás
      </button>
      <button
        @click="() => emits('delete', openHour.id)"
        class="bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-3 rounded transition-colors"
      >
        Törlés
      </button>
    </div>
  </div>
</template>
