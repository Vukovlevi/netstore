<script setup lang="ts">
import { ref } from "vue";
import { ContractTypeClass, type ContractType } from "../../types/ContractType";
import Modal from "../Modal.vue";

const NEW_CONTRACT_TYPE = "Új szerződéstípus felvitele";
const UPDATE_CONTRACT_TYPE = "Szerződéstípus módosítása";
const MAX_WORK_HOURS_A_DAY = 16;
const DAYS_A_WEEK = 7;
const MAX_WEEKLY_WORK_HOURS = DAYS_A_WEEK * MAX_WORK_HOURS_A_DAY;

const props = defineProps<{ contractType: ContractType | null }>();
const emits = defineEmits(["feedback", "back"]);
const contractType = ref(new ContractTypeClass(props.contractType));
let oldContractType = contractType.value.toContractType();
const isModalOpen = ref(false);
const isUpdate = ref(props.contractType != null);

function validate(): { message: string; valid: boolean } {
  if (isUpdate.value && contractType.value.id == 0) {
    return {
      message:
        "A szerződéstípus azonosítója hiányzik, próbáld meg frissíteni az oldalt!",
      valid: false,
    };
  }

  if (contractType.value.name == "" || contractType.value.weeklyHours == 0)
    return {
      message: "A *-gal jelölt mezők kitöltése kötelező!",
      valid: false,
    };

  if (
    contractType.value.weeklyHours <= 0 ||
    contractType.value.weeklyHours > MAX_WEEKLY_WORK_HOURS ||
    contractType.value.weeklyHours !=
      parseInt(contractType.value.weeklyHours.toString())
  )
    return {
      message: `A szerződéstípus heti munkaóráinak száma nem valid (1-${MAX_WEEKLY_WORK_HOURS} közötti egész szám kell legyen!)`,
      valid: false,
    };

  return { message: "", valid: true };
}

async function saveContractType() {
  const { message, valid } = validate();
  if (!valid) {
    emits("feedback", "warning", message, null, false);
    return;
  }

  try {
    const resp = await fetch("/api/contract-type", {
      method: isUpdate.value ? "PUT" : "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contractType.value.toContractType()),
    });
    const data = await resp.json();

    if (data.error) {
      emits(
        "feedback",
        "error",
        "A mentés közben hiba történt: " + data.error,
        null,
        false
      );
      return;
    }

    oldContractType = contractType.value.toContractType();

    emits(
      "feedback",
      "success",
      data.message,
      contractType.value.toContractType(),
      isUpdate.value
    );
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt a következő műveletet nem sikerült végrehatjani: " +
        isUpdate.value
        ? UPDATE_CONTRACT_TYPE
        : NEW_CONTRACT_TYPE,
      null,
      false
    );
  }
}

function cancel() {
  isModalOpen.value = false;
}

function confirm() {
  isModalOpen.value = false;
  emits("back");
}

//TODO: összes helyen törlésnél modal, meg esetleg a mégse gomboknál is, egyeztess difivel
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
          Szerződéstípus adatai
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
          {{ isUpdate ? UPDATE_CONTRACT_TYPE : NEW_CONTRACT_TYPE }}
        </p>
      </div>
    </div>

    <form class="space-y-6" @submit.prevent="saveContractType">
      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="name"
        >
          Név*
        </label>
        <div class="mt-1">
          <input
            id="name"
            type="text"
            placeholder="Adja meg a szerződéstípus nevét"
            v-model="contractType.name"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            required
          />
        </div>
      </div>

      <div>
        <label
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          for="weekly-hours"
        >
          Heti munkaórák*
        </label>
        <div class="mt-1">
          <input
            id="weekly-hours"
            type="number"
            placeholder="Adja meg a heti munkaórák számát"
            v-model="contractType.weeklyHours"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            required
          />
        </div>
      </div>

      <div class="flex justify-end gap-3 pt-4">
        <button
          type="button"
          @click="
            () => {
              if (!contractType.compare(oldContractType)) isModalOpen = true;
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
