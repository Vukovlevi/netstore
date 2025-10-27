import { createWebHistory, createRouter } from "vue-router";

import Home from "../views/Home.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/login", component: () => import("../views/Login.vue") },
  { path: "/users", component: () => import("../views/Users.vue") },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
