<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import {
  OpenHourClass,
  type OpenHour,
  type WeekDay,
} from "../../types/OpenHour";
import Modal from "../Modal.vue";

const NEW_OPEN_HOUR = "Új nyitvatartási idő felvitele";
const UPDATE_OPEN_HOUR = "Nyitvatartási idő módosítása";

const props = defineProps<{ openHour: OpenHour | null }>();
const emits = defineEmits(["feedback", "back"]);
const openHour = ref(new OpenHourClass(props.openHour));
let oldOpenHour = openHour.value.toOpenHour();
const isModalOpen = ref(false);
const weekDays: Ref<WeekDay[], WeekDay[]> = ref([]);
const isUpdate = ref(props.openHour != null);
const isSpecialDate = ref(false);

function validate(): { message: string; valid: boolean } {
  if (isUpdate.value && openHour.value.id == 0) {
    return {
      message:
        "A nyitvatartási idő azonosítója hiányzik, próbáld meg frissíteni az oldalt!",
      valid: false,
    };
  }

  if (
    (openHour.value.opensAt == "" || openHour.value.closesAt == "") &&
    !isSpecialDate.value
  ) {
    return {
      message:
        "A nyitvatartási idő kezdetét és végét kötelező megadni egy általános nap esetén!",
      valid: false,
    };
  } else if (
    (openHour.value.opensAt == "" || openHour.value.closesAt == "") &&
    isSpecialDate.value
  ) {
    openHour.value.opensAt = "00:00";
    openHour.value.closesAt = "00:00";
  }

  if (
    /\d{2}:\d{2}^$/.test(openHour.value.opensAt) ||
    /\d{2}:\d{2}^$/.test(openHour.value.closesAt)
  ) {
    return {
      message:
        "A nyitvatartási idő kezdete vagy vége nem a megfelelő formátumban van!",
      valid: false,
    };
  }

  if (isSpecialDate.value && openHour.value.inputSpecialDate == "") {
    return {
      message: "Különleges dátum esetén a dátum megadása kötelező!",
      valid: false,
    };
  }

  if (!isSpecialDate.value && openHour.value.weekDayIds.length == 0) {
    return {
      message: "A nyitvatartási időnek legalább 1 napon érvényben kell lennie!",
      valid: false,
    };
  }

  return { message: "", valid: true };
}

async function saveOpenHour() {
  const { message, valid } = validate();
  if (!valid) {
    emits("feedback", "warning", message, null, false);
    return;
  }

  if (isSpecialDate.value) {
    openHour.value.weekDayIds = [];
  } else {
    openHour.value.inputSpecialDate = "";
  }

  try {
    const resp = await fetch("/api/open-hour", {
      method: isUpdate.value ? "PUT" : "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(openHour.value.toOpenHour()),
    });
    const data = await resp.json();

    if (data.error) {
      emits("feedback", "error", data.error, null, false);
      return;
    }

    setWeekDaysByIds();
    oldOpenHour = openHour.value.toOpenHour();

    emits(
      "feedback",
      "success",
      data.message,
      openHour.value.toOpenHour(),
      isUpdate.value
    );
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt a következő műveletet nem sikerült végrehajtani: " +
        isUpdate.value
        ? UPDATE_OPEN_HOUR
        : NEW_OPEN_HOUR,
      null,
      false
    );
  }
}

function setWeekDaysByIds() {
  openHour.value.weekDays = [];
  openHour.value.weekDayIds.forEach((id) => {
    const weekDay = weekDays.value.find((x) => x.id == id)!.name;
    openHour.value.weekDays.push(weekDay);
  });
}

async function getWeekDays() {
  try {
    const resp = await fetch("/api/weekdays");
    const data = await resp.json();

    if (data.error) {
      emits("feedback", "error", data.error, null, false);
      return;
    }

    weekDays.value = data as WeekDay[];
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült lekérni a hét napjait!",
      null,
      false
    );
  }
}

function setWeekDayIdsByName() {
  openHour.value.weekDayIds = [];
  openHour.value.weekDays.forEach((day) => {
    const id = weekDays.value.find((x) => x.name == day)!.id;
    openHour.value.weekDayIds.push(id);
  });
}

onMounted(async () => {
  await getWeekDays();
  setWeekDayIdsByName();
  if (openHour.value.specialDate?.Valid) {
    isSpecialDate.value = true;
  }
});

function cancel() {
  isModalOpen.value = false;
}

function confirm() {
  isModalOpen.value = false;
  emits("back");
}
</script>

<template>
  <Modal
    v-if="isModalOpen"
    title="Biztosan vissza akar lépni?"
    message="A szerződéstípus adatainak nem mentett módosításai elvesznek!"
    confirm-text="Igen, visszalépek"
    @cancel="cancel"
    @confirm="confirm"
  />
  <div class="container mx-auto max-w-2xl px-4 py-10 sm:px-6 lg:px-8">
    <div class="space-y-8">
      <div>
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
          Nyitvatartás beállítása
        </h2>
        <p class="mt-2 text-lg text-gray-500 dark:text-gray-400">
          {{ isUpdate ? UPDATE_OPEN_HOUR : NEW_OPEN_HOUR }}
        </p>
      </div>
    </div>

    <form class="space-y-6" @submit.prevent="saveOpenHour">
      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="opening-time"
        >
          Nyitás ideje* (hh:MM)
        </label>
        <div class="mt-1">
          <input
            id="opens-at"
            type="time"
            placeholder="08:00"
            v-model="openHour.opensAt"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
          />
        </div>
      </div>

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="closing-time"
        >
          Zárás ideje* (hh:MM)
        </label>
        <div class="mt-1">
          <input
            id="closes-at"
            type="time"
            v-model="openHour.closesAt"
            placeholder="16:00"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
          />
        </div>
      </div>

      <div>
        <input
          type="checkbox"
          name="is-special-date"
          id="is-special-date"
          v-model="isSpecialDate"
          class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50"
        />
        <label for="is-special-date" class="ml-2"
          >Különleges dátum szerinti nyitvatartás</label
        >
      </div>

      <div v-if="isSpecialDate">
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="effective-from"
        >
          Dátum*
        </label>
        <div class="mt-1">
          <input
            id="special-date"
            type="date"
            v-model="openHour.inputSpecialDate"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
          />
        </div>
      </div>

      <div v-else>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          Napok*
        </label>
        <div class="mt-2 space-y-2">
          <div
            v-for="weekDay in weekDays"
            :key="weekDay.id"
            class="flex items-center gap-2"
          >
            <input
              type="checkbox"
              :id="'day-' + weekDay.id"
              :value="weekDay.id"
              v-model="openHour.weekDayIds"
              @click="() => (openHour.wasWeekDayChanged = true)"
              class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50"
            />
            <label
              :for="'day-' + weekDay.id"
              class="text-sm text-gray-700 dark:text-gray-300"
            >
              {{ weekDay.name }}
            </label>
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-3 pt-4">
        <button
          type="button"
          @click="
            () => {
              if (!openHour.compare(oldOpenHour)) isModalOpen = true;
              else confirm();
            }
          "
          class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
        >
          Mégse
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
</template>
