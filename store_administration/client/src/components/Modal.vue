<script setup lang="ts">
const props = defineProps<{
  title: string;
  message: string;
  confirmText: string;
}>();
const emits = defineEmits(["cancel", "confirm"]);
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
  >
    <div
      class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl dark:bg-background-dark"
    >
      <div class="flex items-start gap-4">
        <div
          class="flex h-10 w-10 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30"
        >
          <svg
            class="h-5 w-5 text-red-600 dark:text-red-400"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 9v4m0 4h.01M10.29 3.86l-7.4 12.82A1 1 0 003.7 18h16.6a1 1 0 00.86-1.32l-7.4-12.82a1 1 0 00-1.72 0z"
            />
          </svg>
        </div>

        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ props.title }}
          </h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ props.message }}
          </p>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <button
          type="button"
          class="rounded bg-background-light px-4 py-2 text-sm font-bold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-background-dark/50 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-background-dark"
          @click="emits('cancel')"
        >
          Mégse
        </button>

        <button
          type="button"
          class="rounded bg-red-600 px-4 py-2 text-sm font-bold text-white shadow-sm hover:bg-red-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
          @click="emits('confirm')"
        >
          {{ props.confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>
