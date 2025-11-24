<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import type { WeekDay } from "../../types/OpenHour";
import type { ContractDay } from "../../types/Contract";

type ScheduleDay = {
  id: number;
  weekDayId: number;
  name: string;
  enabled: boolean;
  startingTime: string;
  endingTime: string;
};

const { weekDays, contractDays } = defineProps<{
  weekDays: WeekDay[];
  contractDays: ContractDay[];
}>();
const schedule: Ref<ScheduleDay[], ScheduleDay[]> = ref([]);
const emits = defineEmits(["back", "save"]);

onMounted(() => {
  weekDays.forEach((weekDay) => {
    const contractDay = contractDays.find(
      (x) => x.weekDayId == weekDay.id || x.weekDay == weekDay.name
    );
    const scheduleDay = {
      id: 0,
      weekDayId: weekDay.id,
      name: weekDay.name,
      enabled: false,
      startingTime: "",
      endingTime: "",
    };

    if (contractDay) {
      scheduleDay.id = contractDay.id;
      scheduleDay.enabled = true;
      scheduleDay.startingTime = contractDay.startingTime;
      scheduleDay.endingTime = contractDay.endingTime;
    }

    schedule.value.push(scheduleDay);
  });
});

function saveContractDays() {
  const days: ContractDay[] = [];
  schedule.value.forEach((day) => {
    if (!day.enabled) return;
    days.push({
      id: day.id,
      weekDayId: day.weekDayId,
      startingTime: day.startingTime,
      endingTime: day.endingTime,
      weekDay: day.name,
      deletedAt: null,
    });
  });

  emits("save", days);
}
</script>

<template>
  <div class="container mx-auto max-w-2xl px-4 py-10 sm:px-6 lg:px-8">
    <div class="space-y-8">
      <div>
        <h4 class="text-3xl font-bold text-gray-900 dark:text-white">
          Munkarend beállítása
        </h4>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
          Adja meg, mely napokon és mettől meddig dolgozik a munkavállaló
        </p>
      </div>

      <form class="space-y-6" @submit.prevent="saveContractDays">
        <div
          v-for="day in schedule"
          :key="day.id"
          class="border rounded p-4 bg-white shadow-sm dark:bg-background-dark/50 dark:border-gray-700"
        >
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ day.name }}
            </h3>

            <div class="flex items-center">
              <input
                type="checkbox"
                class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:focus:ring-primary"
                :id="'day-' + day.weekDayId"
                v-model="day.enabled"
              />
              <label
                class="ml-2 text-sm text-gray-700 dark:text-gray-300"
                :for="'day-' + day.weekDayId"
              >
                Dolgozik ezen a napon
              </label>
            </div>
          </div>

          <div v-if="day.enabled" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Kezdés
              </label>
              <input
                type="time"
                class="mt-1 block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                v-model="day.startingTime"
              />
            </div>

            <div>
              <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Befejezés
              </label>
              <input
                type="time"
                class="mt-1 block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                v-model="day.endingTime"
              />
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3 pt-4">
          <button
            type="button"
            @click="emits('back')"
            class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
          >
            Vissza
          </button>

          <button
            type="submit"
            class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          >
            Mentés
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
