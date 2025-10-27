<script setup lang="ts">
import type { User } from '../../types/User';
const { user } = defineProps<{ user: User }>()
</script>

<template>
    <tr>
        <td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900 dark:text-white">
            {{ user.firstname + " " + user.lastname }}
        </td>
        <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">{{ user.username }}</td>
        <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">{{ user.phoneNumber }}</td>
        <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">
            {{ user.email }}</td>
        <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">
            {{ user.role }}</td>
    </tr>
</template>