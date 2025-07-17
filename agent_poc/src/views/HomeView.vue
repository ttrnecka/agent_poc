<script setup>
import MainView from '../components/MainView.vue'
</script>

<template>
  <main>
    <MainView />
  </main>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
