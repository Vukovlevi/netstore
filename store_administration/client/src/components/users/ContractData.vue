<script setup lang="ts">
import { watchEffect } from "vue";
import type { ContractClass } from "../../types/Contract";
import type { ContractType } from "../../types/ContractType";

const { contract, contractTypes } = defineProps<{
  contract: ContractClass;
  contractTypes: ContractType[];
}>();
const emits = defineEmits(["back", "next", "delete"]);

watchEffect(() => {
  if (contract.contractType != "") {
    const contractType = contractTypes.find(
      (x) => x.name == contract.contractType
    );
    if (contractType) {
      contract.contractTypeId = contractType.id;
    }
  }
});
</script>

<template>
  <form class="space-y-6" @submit.prevent="emits('next')">
    <div>
      <label
        class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        for="contract_type"
      >
        Szerződés típusa*
      </label>
      <div class="mt-1">
        <select
          id="contract_type"
          class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
          v-model="contract.contractTypeId"
        >
          <option value="0">Válasszon típust</option>
          <option
            v-for="contractType in contractTypes"
            :key="contractType.id"
            :value="contractType.id"
          >
            {{ contractType.name + " - " + contractType.weeklyHours }}
          </option>
        </select>
      </div>
    </div>

    <div>
      <label
        for="salary"
        class="block text-sm font-medium text-gray-700 dark:text-gray-300"
      >
        Fizetés (Ft)*
      </label>
      <div class="mt-1">
        <input
          id="salary"
          type="number"
          placeholder="350000"
          class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
          v-model="contract.salary"
        />
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label
          for="starts_at"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          Szerződés kezdete*
        </label>
        <div class="mt-1">
          <input
            id="starts_at"
            type="date"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            v-model="contract.startsAt"
          />
        </div>
      </div>

      <div>
        <label
          for="ends_at"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          Szerződés vége
        </label>
        <div class="mt-1">
          <input
            id="ends_at"
            type="date"
            class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
            v-model="contract.inputEndsAt"
          />
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-3 pt-4">
      <button
        v-if="contract.id != 0"
        @click="() => emits('delete')"
        class="bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-3 rounded"
      >
        Szerződés törlés
      </button>
      <button
        class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
        type="button"
        @click="emits('back')"
      >
        Vissza
      </button>

      <button
        class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        type="submit"
      >
        Következő
      </button>
    </div>
  </form>
</template>
