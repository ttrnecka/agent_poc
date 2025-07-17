<script setup>
import { ref, onMounted, computed } from 'vue'
import { dataStore } from '@/stores/store'
const collectorStatus = ref("UNKNOWN")
const conn = ref(null)
const loadingText = "Loading..."
const dStore = dataStore()
// function send() {
//   if (!conn.value) {
//       return false;
//   }
//   if (!msg.value.value) {
//       return false;
//   }
//   conn.value.send(msg.value.value)
//   msg.value.value = "";
//   return false;
// }

onMounted(() => {
  conn.value = new WebSocket("ws://localhost:8888/ws");

  conn.value.onmessage = function (evt) {
    console.log(evt.data)
    const message = JSON.parse(evt.data)
    dStore.updateCollectorState(message['Source'],{status: message['Text']})
      // collectorStatus.value=evt.data
  };
})

const loadedMessage = computed(() => {
  return dStore.fetchError ? dStore.fetchError.message : loadingText
})
</script>
<template>
<div>
  <p v-if="!dStore.collectors">{{ loadedMessage }}</p>
  <div v-else>
    <table class="table">
      <thead class="thead-dark">
        <tr>
          <th scope="col">Collector Name</th>
          <th scope="col">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(coll, index) in dStore.collectors">
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