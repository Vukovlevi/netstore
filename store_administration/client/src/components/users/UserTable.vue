<script setup lang="ts">
import type { User } from "../../types/User";
import UserRow from "./UserRow.vue";

const { users } = defineProps<{ users: User[] }>();
const emits = defineEmits(["modify", "delete"])
</script>

<template>
  <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
    <thead class="bg-gray-50 dark:bg-gray-700/50">
      <tr>
        <th
          class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
          scope="col"
        >
          Név
        </th>
        <th
          class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
          scope="col"
        >
          Felhasználónév
        </th>
        <th
          class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
          scope="col"
        >
          Telefonszám
        </th>
        <th
          class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
          scope="col"
        >
          Email cím
        </th>
        <th
          class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
          scope="col"
        >
          Rang
        </th>
        <th
          class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
          scope="col"
        >
          Vezérlők
        </th>
      </tr>
    </thead>
    <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
      <UserRow v-for="user in users" :user="user" @modify="(user: User) => emits('modify', user)" @delete="(userId: number) => emits('delete', userId)" />
    </tbody>
  </table>
</template>
