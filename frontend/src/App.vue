<template>
  <div>
    <div class="nav">
      <router-link to="/">Home</router-link>
      <router-link to="/posts">Posts</router-link>
    </div>
    <router-view />
  </div>
</template>
<script>
export default {
  name: "App",
  mounted: function () {
    if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
      document.documentElement.classList.toggle("dark");
    }

    window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
      document.documentElement.classList.toggle("dark");
      // localStorage.setItem("theme", e.matches ? "dark" : "light")
    });
  },
};
</script>
