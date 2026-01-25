<script setup lang="ts">
import { onMounted, ref, watchEffect, type Ref } from "vue";
import {
  ContractClass,
  type Contract,
  type ContractDay as TContractDay,
} from "../../types/Contract.ts";
import type { ContractType } from "../../types/ContractType.ts";
import ContractData from "./ContractData.vue";
import ContractDay from "./ContractDay.vue";
import type { WeekDay } from "../../types/OpenHour.ts";
import Modal from "../Modal.vue";

const NEW_CONTRACT = "Szerződés hozzáadása";
const UPDATE_CONTRACT = "Szerződés módosítása";

const props = defineProps<{ userId: number }>();
const emits = defineEmits(["feedback", "back", "next"]);
const contract = ref(new ContractClass());
let oldContract: Contract;
let oldInputEndsAt = "";
let wasDeleted = false;
let loaded = false;
const isUpdate = ref(false);
const contractTypes: Ref<ContractType[], ContractType[]> = ref([]);
const mode: Ref<"data" | "day", "data" | "day"> = ref("data");
const weekDays: Ref<WeekDay[], WeekDay[]> = ref([]);
const file: Ref<File | null, File | null> = ref(null)
const isModalOpen = ref(false)

async function getContractTypes() {
  try {
    const resp = await fetch("/api/contract-type");
    const data = await resp.json();

    if (data.error) {
      return emits("feedback", "error", data.error, null, false);
    }

    contractTypes.value = data as ContractType[];
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült betölteni a szerződéstípusokat!",
      null,
      false
    );
  }
}

async function getContract() {
  try {
    const resp = await fetch(`/api/contract?userId=${props.userId}`);
    if (resp.status == 204) {
      contract.value.userId = props.userId;
      return;
    }
    const data = await resp.json();

    if (data.error) {
      return emits("feedback", "error", data.error, null, false);
    }

    isUpdate.value = true;
    contract.value = new ContractClass(data as Contract);
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült lekérdezni a szerződést!",
      null,
      false
    );
  }
  contract.value.userId = props.userId;
}

async function getWeekDays() {
  try {
    const resp = await fetch("/api/weekdays");
    const data = await resp.json();

    if (data.error) {
      return emits("feedback", "error", data.error, null, false);
    }

    weekDays.value = data as WeekDay[];
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült lekérdezni a hét napjait!",
      null,
      false
    );
  }
}

function onFileChanged(event: Event) {
  const target = event.target as HTMLInputElement

  if (!target.files || target.files.length == 0) {
    file.value = null
    return
  }

  file.value = target.files[0]!
}

async function saveContract() {
  const { message, valid } = validate();
  if (!valid) {
    emits("feedback", "warning", message, null, false);
    return;
  }

  try {
    const formData = new FormData()
    formData.append("contract", JSON.stringify(contract.value.toContract()))
    if (file.value) formData.append("file", file.value)
    const resp = await fetch("/api/contract", {
      method: isUpdate.value ? "PUT" : "POST",
      headers: { "Content-Type": "application/json" },
      body: formData,
    });
    const data = await resp.json();

    if (data.error) {
      return emits("feedback", "error", data.error, null, false);
    }

    contract.value.contractType = contractTypes.value.find(contractType => contractType.id == contract.value.contractTypeId)!.name
    emits("feedback", "success", data.message, null, false);
    oldContract = contract.value.toContract()
    contract.value.changedEndsAt = false
    contract.value.changedContractDays = false
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült menteni a szerződést!",
      null,
      false
    );
  }
}

function validate(): { message: string; valid: boolean } {
  if (isUpdate.value) return validateUpdate();
  else return validateInsert();
}

function validateUpdate(): { message: string; valid: boolean } {
  if (contract.value.id == 0) {
    return {
      message:
        "A szerződés azonosítója hiányzik! Próbálja újratölteni az oldalt!",
      valid: false,
    };
  }

  return validateInsert();
}

function validateInsert(): { message: string; valid: boolean } {
  if (contract.value.userId == 0) {
    return {
      message:
        "A felhasználó azonosítója nem található! (Próbálja frissíteni az oldalt!)",
      valid: false,
    };
  }

  if (
    contract.value.contractTypeId == 0 ||
    contract.value.salary == 0 ||
    contract.value.startsAt == ""
  ) {
    return {
      message:
        "A *-gal jelölt mezők kitöltése kötelező az általános adatoknál!",
      valid: false,
    };
  }

  if (contract.value.contractDays.length == 0) {
    return {
      message: "Legalább 1 munkanap adatait fel kell vinni!",
      valid: false,
    };
  }

  for (const day of contract.value.contractDays) {
    if (day.startingTime == "" || day.endingTime == "") {
      return {
        message: "Egy munkanapnak kötelező kezdési és vég időpontot állítani!",
        valid: false,
      };
    }
  }

  //todo: check against contract type weekly hours the total working hours of a user

  return { message: "", valid: true };
}

async function deleteContract() {
  if (contract.value.id == 0) {
    emits(
      "feedback",
      "warning",
      "A szerződés törléséhez szükséges az azonosítója! (Próbálja frissíteni az oldalt!)",
      null,
      false
    );
  }

  try {
    const resp = await fetch("/api/contract", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contract.value.toContract()),
    });

    if (resp.status == 204) {
      emits("feedback", "success", "Szerződés törlése sikeres!", null, false);
      wasDeleted = true;
      return;
    }

    const data = await resp.json();
    if (data.error) {
      emits("feedback", "error", data.error, null, false);
    }
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült törölni a szerződést!",
      null,
      false
    );
  }
}

onMounted(async () => {
  getWeekDays();
  await getContractTypes();
  await getContract();
  if (contract.value.contractType != "") contract.value.contractTypeId = contractTypes.value.find(contractType => contractType.name == contract.value.contractType)!.id
  oldContract = contract.value.toContract()
});

watchEffect(() => {
  console.log(contract.value.inputEndsAt)
  if (!loaded) {
    loaded = true
    return
  }
  if (oldInputEndsAt != contract.value.inputEndsAt) {
    oldInputEndsAt = contract.value.inputEndsAt
    contract.value.changedEndsAt = true
  }
})

async function deleteContractFile() {
  try {
    const resp = await fetch("/api/contract-file", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contract.value.toContract())
    })

    if (resp.status == 204) {
      emits("feedback", "success", "Szerződés fájl törlése sikeres!", null, false);
      return
    }

    const data = await resp.json();
    if (data.error) {
      emits("feedback", "error", data.error, null, false);
    }
  } catch (err) {
    console.error(err);
    emits(
      "feedback",
      "error",
      "Ismeretlen hiba miatt nem sikerült törölni a szerződés fájlt!",
      null,
      false
    );
  }
}

function cancel() {
  isModalOpen.value = false
}

function confirm() {
  isModalOpen.value = false
  deleteContractFile()
}
</script>

<template>
  <Modal v-if="isModalOpen" title="Megerősítés" message="Biztosan törölni akarja a szerződés fájlt?"
    confirm-text="Törlés" @cancel="cancel" @confirm="confirm" />
  <div class="container mx-auto max-w-2xl px-4 py-10 sm:px-6 lg:px-8 max-h-[calc(100vh-15rem)] overflow-y-auto">
    <div class="space-y-8">
      <div>
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
          Szerződés adatai
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
          {{ isUpdate ? UPDATE_CONTRACT : NEW_CONTRACT }}
        </p>
      </div>

      <ContractData v-if="mode == 'data'" :contract="contract" :contractTypes="contractTypes" :file="file"
        @fileChanged="onFileChanged" @clearFile="() => file = null"
        @back="emits('back', contract.compare(oldContract) || wasDeleted)" @next="mode = 'day'"
        @delete="deleteContract" />
      <ContractDay v-else
        @back="(contractDays: TContractDay[], changedContractDays: boolean) => { contract.contractDays = contractDays; if (!contract.changedContractDays) contract.changedContractDays = changedContractDays; mode = 'data' }"
        :weekDays="weekDays" :contractDays="contract.contractDays"
        @save="(contractDays: TContractDay[]) => { contract.contractDays = contractDays; saveContract() }" />
    </div>
  </div>
</template>
