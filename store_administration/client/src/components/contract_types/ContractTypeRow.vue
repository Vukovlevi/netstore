<script setup lang="ts">
import type { ContractType } from '../../types/ContractType';

const { contractType } = defineProps<{ contractType: ContractType }>();
const emits = defineEmits(["modify", "delete"])
</script>

<template>
    <tr>
        <td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900 dark:text-white">
            {{ contractType.name }}
        </td>
        <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">
            {{ contractType.weeklyHours }}
        </td>
        <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-gray-300">
            <div class="flex flex-col space-y-2">
                <button @click="() => { emits('modify', contractType) }"
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-3 rounded">
                    Módosítás
                </button>
                <button @click="() => emits('delete', contractType.id)"
                    class="bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-3 rounded">
                    Törlés
                </button>
            </div>
        </td>
    </tr>
</template>
