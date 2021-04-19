<template>
  <div class="center-vh text-center">
    <div>
      <img class="rounded" src="@/assets/images/profile.jpg" width="256" height="256" alt="" />
      <h1 class="title">Taliesin Millhouse</h1>
      <h2 class="sub-title">Software Developer</h2>
      <a class="icon" href="https://github.com/strongishllama" target="_blank" rel="noopener" aria-label="GitHub">
        <i class="fab fa-github fa-2x"></i>
      </a>
      <a class="icon" href="https://www.linkedin.com/in/taliesin-millhouse-00001" target="_blank" rel="noopener" aria-label="LinkedIn">
        <i class="fab fa-linkedin-in fa-2x"></i>
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: "Home",
};
</script>
<style lang="scss" scoped>
.center-vh {
  height: 80vh;
}
</style>
