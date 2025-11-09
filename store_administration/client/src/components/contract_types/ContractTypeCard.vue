<script setup lang="ts">
import type { ContractType } from '../../types/ContractType';

const { contractType } = defineProps<{ contractType: ContractType }>();
const emits = defineEmits(["modify", "delete"])
</script>

<template>
    <div class="bg-white shadow-md rounded-lg p-5 max-w-sm hover:shadow-lg transition-shadow">
        <h3 class="text-lg font-semibold mb-2">{{ contractType.name }}</h3>
        <p class="text-gray-600 text-sm mb-4">
            {{ "Heti munkaórák: " + contractType.weeklyHours }}
        </p>

        <div class="flex gap-2">
            <button @click="() => emits('modify', contractType)"
                class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-3 rounded transition-colors">
                Módosítás
            </button>
            <button @click="() => emits('delete', contractType.id)"
                class="bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-3 rounded transition-colors">
                Törlés
            </button>
        </div>
    </div>
</template>
