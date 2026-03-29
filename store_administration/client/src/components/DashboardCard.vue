<script setup lang="ts">
defineProps({
    title: String,
    description: String,
    external: Boolean,
});
</script>

<template>
    <div @click="$emit('click')"
        class="cursor-pointer rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-background-dark/50 p-6 shadow-sm hover:shadow-md transition-all hover:-translate-y-1 hover:border-primary">
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <slot name="icon">
                    <!-- Default icon -->
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12h6m-6 4h6M9 8h6M5 6h14M5 18h14" />
                    </svg>
                </slot>
            </div>

            <!-- Content -->
            <div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                    {{ title }}
                </h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    {{ description }}
                </p>
            </div>
        </div>

        <!-- External badge -->
        <div v-if="external" class="mt-4 text-xs text-primary font-medium">
            Külső oldal →
        </div>
    </div>
</template>