<script setup lang="ts">
import Header from "./components/Header.vue";
import { RouterView } from "vue-router";
</script>

<template>
  <Header class="absolute"></Header>
  <main class="p-3">
    <RouterView />
  </main>
</template>

<style scoped></style>
