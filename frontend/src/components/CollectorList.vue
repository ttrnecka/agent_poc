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
<div class="container-fluid">
  <p v-if="!apiStore.sortedCollectors">{{ loadedMessage }}</p>
  <div v-else class="row">
    <div class="col-auto" style="flex: 0 0 200px;">
    </div>
    <div class="col">
      <table class="table">
        <thead class="thead-dark">
          <tr>
            <th scope="col">Collector Name</th>
            <th scope="col">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(coll, index) in apiStore.collectors" :key="index">
            <td>{{coll.key}}</td>
            <td>{{coll.status || collectorStatus}}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</div>
 
</template>

<style type="text/css">

</style>