import { createWebHistory, createRouter } from "vue-router";

import Home from "../views/Home.vue";

const routes = [
  { path: "/", component: Home, meta: { requiresAuth: true } },
  { path: "/login", component: () => import("../views/Login.vue") },
  {
    path: "/users",
    component: () => import("../views/Users.vue"),
    meta: { requiresAuth: true },
  },
  {path: "/password-change", component: () => import("../components/profile/PasswordChange.vue"), meta: {requiresAuth: true},},
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, _, next) => {
  const authState = await isAuthenticated();
  if (to.meta.requiresAuth && !authState) {
    next("/login");
  } else {
    next();
  }
});

export { router };

async function isAuthenticated(): Promise<boolean> {
  try {
    const resp = await fetch("/api/echo");
    return !resp.redirected;
  } catch (err) {
    console.error(err);
    return false;
  }
}
