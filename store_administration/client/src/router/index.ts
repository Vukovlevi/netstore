import { createWebHistory, createRouter } from "vue-router";

import Home from "../views/Home.vue";

const PASSWORD_CHANGE_URL = "password-change";

const routes = [
  { path: "/", component: Home, meta: { requiresAuth: true } },
  { path: "/login", component: () => import("../views/Login.vue") },
  {
    path: "/users",
    component: () => import("../views/Users.vue"),
    meta: { requiresAuth: true },
  },
  {path: "/password-change", component: () => import("../components/profile/PasswordChange.vue"),},
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const authState = await isAuthenticated();

  if (authState.passwordChangeUrl != "") {
    console.log("ide belépünk")
    next("/" + authState.passwordChangeUrl);
    return;
  }

  if (to.meta.requiresAuth && !authState.valid) {
    next("/login");
  } else {
    next();
  }
});

export { router };

async function isAuthenticated(): Promise<{valid: boolean, passwordChangeUrl: string}> {
  try {
    const resp = await fetch("/api/echo");

    const passChangeUrlParts = resp.url.split("/");
    let lastPart = passChangeUrlParts[passChangeUrlParts.length - 1] || "";
    if (!resp.redirected || lastPart != PASSWORD_CHANGE_URL) {
      lastPart = "";
    }

    return {valid: !resp.redirected, passwordChangeUrl: lastPart};
  } catch (err) {
    console.error(err);
    return {valid: false, passwordChangeUrl: ""};
  }
}
