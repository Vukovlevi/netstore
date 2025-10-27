<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>

<template>
    <header
        class="absolute top-0 left-0 right-0 z-10 flex h-16 items-center border-b border-background-light-darken dark:border-background-dark-lighten bg-background-light dark:bg-background-dark px-4 sm:px-6 lg:px-8">

        <!-- Left: Logo + Title -->
        <div class="flex items-center gap-4">
            <div class="text-primary">
                <svg class="h-6 w-6" fill="none" height="24" stroke="currentColor" stroke-linecap="round"
                    stroke-linejoin="round" stroke-width="2" viewBox="0 0 24 24" width="24"
                    xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L2 7v10l10 5 10-5V7L12 2z"></path>
                    <path d="M2 17l10 5"></path>
                    <path d="M22 17l-10 5"></path>
                    <path d="M12 12l10-5"></path>
                    <path d="M12 12v10"></path>
                    <path d="M12 12L2 7"></path>
                </svg>
            </div>
            <h1 class="text-lg font-bold text-gray-900 dark:text-white">NetStore</h1>
        </div>

        <!-- Center: Nav -->
        <nav
            class="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 hidden items-center gap-6 text-sm font-medium md:flex">
            <RouterLink class="text-gray-600 dark:text-gray-300 hover:text-primary dark:hover:text-primary" to="/">
                Home
            </RouterLink>
            <RouterLink class="text-gray-600 dark:text-gray-300 hover:text-primary dark:hover:text-primary" to="/login">
                Login
            </RouterLink>
            <RouterLink class="text-gray-600 dark:text-gray-300 hover:text-primary dark:hover:text-primary" to="/users">
                Users
            </RouterLink>
        </nav>

        <!-- Right: Profile -->
        <div class="flex items-center gap-4 ml-auto">
            <svg xmlns="http://www.w3.org/2000/svg"
                class="aspect-square w-8 h-8 rounded-full bg-cover bg-center profile-icon" viewBox="0 0 640 640">
                <path
                    d="M320 312C386.3 312 440 258.3 440 192C440 125.7 386.3 72 320 72C253.7 72 200 125.7 200 192C200 258.3 253.7 312 320 312zM290.3 368C191.8 368 112 447.8 112 546.3C112 562.7 125.3 576 141.7 576L498.3 576C514.7 576 528 562.7 528 546.3C528 447.8 448.2 368 349.7 368L290.3 368z" />
            </svg>
        </div>
    </header>

</template>

<style scoped>
.profile-icon {
    fill: rgb(17 115 212 / var(--tw-text-opacity, 1))
}
</style>