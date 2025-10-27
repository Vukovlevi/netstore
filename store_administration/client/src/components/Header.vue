<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>

<template>
    <header>
        <nav>
            <RouterLink to="/">Home</RouterLink>
            <RouterLink to="/login">Login</RouterLink>
        </nav>
    </header>
</template>