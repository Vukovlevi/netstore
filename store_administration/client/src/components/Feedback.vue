<script setup lang="ts">
import { ref } from 'vue';
import type { Feedback } from '../types/Feedback';

const props = defineProps<{feedback: Feedback}>()
const style = ref("")

switch (props.feedback.type) {
    case "info":
        style.value = "p-3 text-sm rounded-lg border border-blue-400 bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-800 mb-3";
        break;
    case "success":
        style.value = "p-3 text-sm rounded-lg border border-green-400 bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-300 dark:border-green-800 mb-3";
        break;
    case "warning":
        style.value = "p-3 text-sm rounded-lg border border-amber-400 bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-800 mb-3";
        break;
    case "error":
        style.value = "p-3 text-sm rounded-lg border border-red-400 bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800 mb-3";
        break;
    default:
        break;
}
</script>

<template>
    <div
      :class="style"
      role="alert"
    >
      {{ props.feedback.message }}
    </div>
</template>