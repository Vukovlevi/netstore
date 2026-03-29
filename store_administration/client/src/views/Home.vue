<script setup lang="ts">
import { useRouter } from "vue-router";
import DashboardCard from "../components/DashboardCard.vue";

const router = useRouter();

function goTo(path: string) {
    router.push(path);
};

function openExternal() {
    window.open(import.meta.env.VITE_PROD_ADMIN_ORIGIN, "_self");
};
</script>

<template>
    <div class="container mt-[5rem] mx-auto max-w-6xl px-4 py-10 sm:px-6 lg:px-8">
        <!-- Header -->
        <div class="mb-10">
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                Főoldal
            </h1>
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
                Válasszon egy menüpontot az alábbiak közül
            </p>
        </div>

        <!-- Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            <DashboardCard title="Felhasználók kezelése"
                description="Felhasználók létrehozása, módosítása, törlése és szerződésük kezelése"
                @click="goTo('/users')" />

            <DashboardCard title="Üzlet adatai"
                description="Üzlet címének, típusának, valamint a központi szerverre csatlakozás adatainak kezelése"
                @click="goTo('/store-detail')" />

            <DashboardCard title="Szerződéstípusok" description="Szerződéstípusok kezelése"
                @click="goTo('/contract-types')" />

            <DashboardCard title="Nyitvatartás" description="Nyitvatartási idő beállítása"
                @click="goTo('/open-hour')" />

            <DashboardCard title="Profil" description="Saját adatok módosítása" @click="goTo('/profile')" />

            <DashboardCard title="Termék kezelő" description="A termék kezelő felület megnyitása" external
                @click="openExternal" />
        </div>
    </div>
</template>