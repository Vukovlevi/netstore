<script setup lang="ts">
import type { ContractType } from '../../types/ContractType';
import ContractTypeRow from './ContractTypeRow.vue';

const { contractTypes } = defineProps<{ contractTypes: ContractType[] }>();
const emits = defineEmits(["modify", "delete"])
</script>

<template>
    <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-700/50">
            <tr>
                <th class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    scope="col">
                    Név
                </th>
                <th class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    scope="col">
                    Heti munkaórák
                </th>
                <th class="px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    scope="col">
                    Vezérlők
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <ContractTypeRow v-for="contractType in contractTypes" :contractType="contractType"
                @modify="(contractType: ContractType) => emits('modify', contractType)"
                @delete="(contractTypeId: number) => emits('delete', contractTypeId)" />
        </tbody>
    </table>
</template>
