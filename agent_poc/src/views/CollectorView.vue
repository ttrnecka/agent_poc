<script setup>
import CollectorList from '../components/CollectorList.vue'
</script>

<template>
  <main>
    <CollectorList />
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
