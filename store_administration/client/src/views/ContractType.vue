<script setup lang="ts">
import { onMounted, ref, type Ref } from 'vue';
import type { FeedbackType, Feedback as TFeedback } from '../types/Feedback';
import type { ContractType } from '../types/ContractType';
import ContractTypeData from '../components/contract_types/ContractTypeData.vue';
import Feedback from '../components/Feedback.vue';
import SearchBar from '../components/SearchBar.vue';
import ContractTypeCard from '../components/contract_types/ContractTypeCard.vue';

const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);
let contractTypes: ContractType[] = [];
const filteredContractTypes: Ref<ContractType[], ContractType[]> = ref([]);
const currentContractType: Ref<ContractType | null, ContractType | null> = ref(null);
const mode: Ref<"all" | "single", "all" | "single"> = ref("all");

async function getContractTypes() {
    try {
        const resp = await fetch("/api/contract-type");
        const data = await resp.json();

        if (data.error) {
            feedback.value = { type: "error", message: (data.error as string) };
            return;
        }

        feedback.value = null;
        contractTypes = data as ContractType[];
        filteredContractTypes.value = contractTypes;
    } catch (err) {
        feedback.value = { type: "error", message: "Ismeretlen hiba miatt nem sikerült lekérni a szerződéstípusokat!" };
        console.error(err);
    }
}

function modifyContractType(contractType: ContractType) {
    currentContractType.value = contractType
    feedback.value = null;
    mode.value = "single"
}

function handleFeedback(type: FeedbackType, msg: string, contractType: ContractType | null, isUpdate: boolean) {
    feedback.value = { type: type, message: msg }
    if (contractType == null) return;

    if (isUpdate) updateContractType(contractType);
    else createNewContractType(contractType);
}

function createNewContractType(contractType: ContractType) {
    contractTypes.push(contractType);
    filteredContractTypes.value = contractTypes;
}

function updateContractType(contractType: ContractType) {
    const idx = contractTypes.findIndex(x => x.id == contractType.id);
    contractTypes[idx] = contractType;
    filteredContractTypes.value = contractTypes;
}

async function deleteContractType(contractTypeId: Number) {
    if (contractTypeId == 0) {
        feedback.value = { type: "warning", message: "A szerződéstípus azonosítója hiányzik, próbáld meg frissíteni az oldalt!" }
        return;
    }

    try {
        const resp = await fetch("/api/contract-type", {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: contractTypeId }),
        })

        if (resp.status == 204) {
            contractTypes = contractTypes.filter(x => x.id != contractTypeId);
            filteredContractTypes.value = contractTypes;
            feedback.value = { type: "success", message: "Szerződéstípus törlése sikeres!" };
            return;
        }

        const data = await resp.json()
        if (data.error) {
            feedback.value = { type: "error", message: data.error };
            return;
        }

        feedback.value = { type: "error", message: "Szerződéstípus törlése sikertelen!" };
    } catch (err) {
        feedback.value = { type: "error", message: "Ismeretlen hiba miatt nem sikerült törölni a szerződéstípust!" };
        console.error(err);
    }
}

function search(searchValue: string) {
    if (searchValue == "") {
        filteredContractTypes.value = contractTypes;
        return;
    }

    filteredContractTypes.value = contractTypes.filter((x) =>
        x.name.toLowerCase().includes(searchValue.toLowerCase()) || x.weeklyHours.toString() == searchValue
    );
}

onMounted(() => {
    getContractTypes();
})
</script>

<template>
    <div class="mx-auto max-w-7xl mt-[5rem]">
        <div class="mb-6 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                Szerződéstípusok
            </h2>
            <button v-if="mode == 'all'"
                class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90"
                @click="() => { mode = 'single'; feedback = null; }">
                Szerződéstípus felvitele
            </button>
        </div>
        <SearchBar search-item="Szerződéstípusok" @search="search" v-if="mode == 'all'" />
        <Feedback v-if="feedback != null" :feedback="feedback" />

        <div v-if="mode == 'all'" class="flex justify-center flex-wrap gap-6 items-center">
            <ContractTypeCard v-for="contractType in filteredContractTypes" :contractType="contractType"
                @modify="(contractType: ContractType) => modifyContractType(contractType)"
                @delete="deleteContractType" />
        </div>
        <ContractTypeData v-else :contractType="currentContractType" @feedback="handleFeedback"
            @back="() => { mode = 'all'; feedback = null; currentContractType = null; }" />
    </div>
</template>
