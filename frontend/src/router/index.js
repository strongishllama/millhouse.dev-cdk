import { createRouter, createWebHistory } from "vue-router";
import Home from "@/views/Home.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/posts",
    name: "Posts",
    component: () => import("@/views/Posts.vue"),
  },
  {
    path: "/posts/graceful-shutdowns-in-golang-with-signal-notify-context",
    name: "GracefulShutdownsInGolangWithSignalNotifyContext",
    component: () => import("@/views/posts/GracefulShutdownsInGolangWithSignalNotifyContext.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.afterEach((to) => {
  switch (to.path) {
    case "/":
      document.title = "Taliesin Millhouse";
      break;
    case "/posts":
      document.title = "Posts · Taliesin Millhouse";
      break;
    case "/posts/graceful-shutdowns-in-golang-with-signal-notify-context":
      document.title = "Graceful Shutdowns in Golang with signal.NotifyContext · Taliesin Millhouse";
      break;
    default:
      document.title = "Taliesin Millhouse";
      break;
  }
});

export default router;
