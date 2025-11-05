<script setup lang="ts">
import type { User } from "../../types/User";
const { user } = defineProps<{ user: User }>();
const emits = defineEmits(["modify"])
</script>

<template>
  <tr>
    <td
      class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900 dark:text-white"
    >
      {{ user.lastname + " " + user.firstname }}
    </td>
    <td
      class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300"
    >
      {{ user.username }}
    </td>
    <td
      class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300"
    >
      {{ user.phoneNumber }}
    </td>
    <td
      class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300"
    >
      {{ user.email }}
    </td>
    <td
      class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300"
    >
      {{ user.role }}
    </td>
    <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">
      <div class="flex flex-col space-y-2">
        <button
          @click="() => {emits('modify', user)}"
          class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-3 rounded"
        >
          Módosítás
        </button>
        <button
          @click=""
          class="bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-3 rounded"
        >
          Törlés
        </button>
      </div>
    </td>
  </tr>
</template>
