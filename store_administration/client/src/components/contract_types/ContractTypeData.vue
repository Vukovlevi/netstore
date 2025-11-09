<script setup lang="ts">
import { ref, type Ref } from "vue";
import { ContractTypeClass, type ContractType } from "../../types/ContractType";
import type { Feedback as TFeedback } from "../../types/Feedback";
import Feedback from "../Feedback.vue";
import { router } from "../../router";

const NEW_CONTRACT_TYPE = "Új szerződéstípus felvitele";
const UPDATE_CONTRACT_TYPE = "Szerződéstípus módosítása";
const MAX_WORK_HOURS_A_DAY = 16;
const DAYS_A_WEEK = 7;
const MAX_WEEKLY_WORK_HOURS = DAYS_A_WEEK * MAX_WORK_HOURS_A_DAY;

const props = defineProps<{ contractType: ContractType | null }>();
const emits = defineEmits(["feedback", "back"]);
const contractType = ref(new ContractTypeClass(props.contractType));
const isUpdate = ref(props.contractType != null);
const feedback: Ref<TFeedback | null, TFeedback | null> = ref(null);

function validate(): { message: string; valid: boolean } {
    if (isUpdate.value && contractType.value.id == 0) {
        return {
            message: "A szerződéstípus azonosítója hiányzik, próbáld meg frissíteni az oldalt!",
            valid: false,
        };
    }

    if (
        contractType.value.name == "" ||
        contractType.value.weeklyHours == 0
    )
        return {
            message: "A *-gal jelölt mezők kitöltése kötelező!",
            valid: false,
        };

    if (contractType.value.weeklyHours <= 0 || contractType.value.weeklyHours > MAX_WEEKLY_WORK_HOURS || contractType.value.weeklyHours != parseInt(contractType.value.weeklyHours.toString()))
        return { message: `A szerződéstípus heti munkaóráinak száma nem valid (1-${MAX_WEEKLY_WORK_HOURS} közötti egész szám kell legyen!)`, valid: false };

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
            emits("feedback", "error", "A mentés közben hiba történt: " + data.error, null, false);
            return;
        }

        emits("feedback", "success", data.message, contractType.value.toContractType(), isUpdate.value);
    } catch (err) {
        console.error(err);
        emits(
            "feedback",
            "error",
            "Ismeretlen hiba miatt a következő műveletet nem sikerült végrehatjani: " +
                isUpdate.value ? UPDATE_CONTRACT_TYPE : NEW_CONTRACT_TYPE,
            null, false
        );
    }
}

//TODO: vissza gomb lenyomásakor figyelmeztetés (modal) -> majd törlésnél is lehet ezt használni (hasonló mint a feedback, csak modal)
</script>
<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-background-dark">
        <form class="w-full max-w-md space-y-6 p-6 rounded-lg shadow-md bg-white dark:bg-background-dark/70"
            @submit.prevent="saveContractType">
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

            <Feedback v-if="feedback != null" :feedback="feedback" />

            <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="name">
                    Név*
                </label>
                <div class="mt-1">
                    <input id="name" type="text" placeholder="Adja meg a szerződéstípus nevét"
                        v-model="contractType.name"
                        class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                        required />
                </div>
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="weekly-hours">
                    Heti munkaórák*
                </label>
                <div class="mt-1">
                    <input id="weekly-hours" type="number" placeholder="Adja meg a heti munkaórák számát"
                        v-model="contractType.weeklyHours"
                        class="block w-full rounded border-gray-300 bg-white shadow-sm focus:border-primary focus:ring-primary dark:border-gray-600 dark:bg-background-dark/50 dark:text-white dark:placeholder-gray-400 dark:focus:border-primary dark:focus:ring-primary"
                        required />
                </div>
            </div>


            <div class="flex justify-end gap-3 pt-4">
                <button type="button" @click="() => router.push('/')"
                    class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark">
                    Mégse
                </button>
                <button type="submit"
                    class="rounded bg-primary px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary">
                    Mentés
                </button>
            </div>
        </form>
    </div>
</template>
