<script setup>
import { ref, onMounted, computed } from 'vue'
import { useApiStore } from '@/stores/apiStore'
const collectorStatus = ref("UNKNOWN")

const loadingText = "Loading..."
const apiStore = useApiStore()

const loadedMessage = computed(() => {
  return apiStore.fetchError ? apiStore.fetchError.message : loadingText
})
</script>
<template>
<div>
  <p v-if="!apiStore.collectors">{{ loadedMessage }}</p>
  <div v-else>
    <table class="table">
      <thead class="thead-dark">
        <tr>
          <th scope="col">Collector Name</th>
          <th scope="col">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(coll, index) in apiStore.collectors" :key="index">
          <td>{{index}}</td>
          <td>{{coll.status || collectorStatus}}</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
 
</template>

<style type="text/css">

</style>