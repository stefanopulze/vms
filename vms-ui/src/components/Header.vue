<script setup lang="ts">

import Logo from "@/components/icons/Logo.vue";
</script>

<template>
  <header id="main-header">
    <div class="flex">
      <Logo class="logo-box"/>
      <h1>VMS Core</h1>
    </div>
    <nav>
      <router-link to="/">Home</router-link>
      <router-link to="/info">Info</router-link>
    </nav>
  </header>
</template>

<style>
#main-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-inline: 1rem;
  border-bottom: 1px solid var(--border-color);

  .logo-box {
    width: 24px;
    height: 24px;
    margin-inline-end: .5rem;
  }

  h1 {
    font-size: 1.1rem;
    margin-block: .8rem;
  }

  nav {
    display: flex;
    gap: 1rem;
  }
}
</style>
