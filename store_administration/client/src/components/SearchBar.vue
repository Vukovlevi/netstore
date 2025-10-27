<script setup lang="ts">
import { ref } from "vue"

const { searchItem } = defineProps<{ searchItem: string }>()
const searchValue = ref("")
const emits = defineEmits(["search"])
//TODO: make enter hit work (button submit, preventdefault)
</script>

<template>
    <div class="mb-6 flex flex-col gap-4 sm:flex-row">
        <div class="relative flex-1">
            <input type="text" placeholder="Search employees..."
                class="w-full rounded-lg border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 py-2 pl-4 pr-10 text-gray-900 dark:text-white placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:border-primary focus:ring-primary"
                v-model="searchValue" name="search" />
            <svg xmlns="http://www.w3.org/2000/svg"
                class="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400 dark:text-gray-500"
                viewBox="0 0 640 640">
                <path
                    d="M480 272C480 317.9 465.1 360.3 440 394.7L566.6 521.4C579.1 533.9 579.1 554.2 566.6 566.7C554.1 579.2 533.8 579.2 521.3 566.7L394.7 440C360.3 465.1 317.9 480 272 480C157.1 480 64 386.9 64 272C64 157.1 157.1 64 272 64C386.9 64 480 157.1 480 272zM272 416C351.5 416 416 351.5 416 272C416 192.5 351.5 128 272 128C192.5 128 128 192.5 128 272C128 351.5 192.5 416 272 416z" />
            </svg>
        </div>

        <button
            class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90"
            @click="emits('search', searchValue)">
            {{ "Search " + searchItem }}
        </button>
    </div>
</template>