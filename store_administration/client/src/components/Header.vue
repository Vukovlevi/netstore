<script setup lang="ts">
import { ref } from "vue";
import { RouterLink } from "vue-router";

const isOpen = ref(false);

const toggleMenu = () => {
  isOpen.value = !isOpen.value;
};
</script>

<template>
  <header
    class="absolute top-0 left-0 right-0 z-10 flex h-16 items-center border-b border-background-light-darken dark:border-background-dark-lighten bg-background-light dark:bg-background-dark px-4 sm:px-6 lg:px-8">
    <!-- Logo -->
    <RouterLink to="/" class="flex items-center gap-4">
      <div class="text-primary">
        <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path d="M12 2L2 7v10l10 5 10-5V7L12 2z"></path>
          <path d="M2 17l10 5"></path>
          <path d="M22 17l-10 5"></path>
          <path d="M12 12l10-5"></path>
          <path d="M12 12v10"></path>
          <path d="M12 12L2 7"></path>
        </svg>
      </div>
      <h1 class="text-lg font-bold text-gray-900 dark:text-white">NetStore</h1>
    </RouterLink>

    <!-- Desktop Nav -->
    <nav
      class="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 hidden lg:flex items-center gap-6 text-sm font-medium">
      <RouterLink class="nav-link" to="/">Főoldal</RouterLink>
      <RouterLink class="nav-link" to="/users">Felhasználók kezelése</RouterLink>
      <RouterLink class="nav-link" to="/store-detail">Üzlet adatai</RouterLink>
      <RouterLink class="nav-link" to="/contract-types">Szerződéstípusok kezelése</RouterLink>
      <RouterLink class="nav-link" to="/open-hour">Nyitvatartási idő beállítása</RouterLink>
    </nav>

    <!-- Right side -->
    <div class="flex items-center gap-4 ml-auto">
      <!-- Hamburger -->
      <button class="lg:hidden text-gray-700 dark:text-gray-300" @click="toggleMenu">
        <svg v-if="!isOpen" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" stroke="currentColor"
          viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" stroke="currentColor"
          viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- Profile -->
      <RouterLink to="/profile" class="flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 rounded-full profile-icon" viewBox="0 0 640 640">
          <path
            d="M320 312C386.3 312 440 258.3 440 192C440 125.7 386.3 72 320 72C253.7 72 200 125.7 200 192C200 258.3 253.7 312 320 312zM290.3 368C191.8 368 112 447.8 112 546.3C112 562.7 125.3 576 141.7 576L498.3 576C514.7 576 528 562.7 528 546.3C528 447.8 448.2 368 349.7 368L290.3 368z" />
        </svg>
      </RouterLink>
    </div>

    <!-- Mobile Menu -->
    <div v-if="isOpen"
      class="absolute top-16 left-0 w-full bg-background-light dark:bg-background-dark border-b border-gray-200 dark:border-gray-700 lg:hidden flex flex-col items-center gap-4 py-4">
      <RouterLink class="nav-link" to="/" @click="toggleMenu">Főoldal</RouterLink>
      <RouterLink class="nav-link" to="/users" @click="toggleMenu">Felhasználók kezelése</RouterLink>
      <RouterLink class="nav-link" to="/store-detail" @click="toggleMenu">Üzlet adatai</RouterLink>
      <RouterLink class="nav-link" to="/contract-types" @click="toggleMenu">Szerződéstípusok kezelése</RouterLink>
      <RouterLink class="nav-link" to="/open-hour" @click="toggleMenu">Nyitvatartási idő beállítása</RouterLink>
    </div>
  </header>
</template>

<style scoped>
.profile-icon {
  fill: rgb(17 115 212 / var(--tw-text-opacity, 1));
}

.nav-link {
  @apply text-gray-600 dark:text-gray-300 hover:text-primary dark:hover:text-primary;
}
</style>
