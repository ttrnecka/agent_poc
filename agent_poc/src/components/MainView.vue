<script setup>
import { ref, onMounted, computed } from 'vue'
import { dataStore } from '@/stores/store'
import { Modal } from "bootstrap";

const newProbe = {
    id: null,
    policy: null,
    collector: null,
    version: null,
    address: null,
    port: null,
    user: null,
    password: null
}
const dStore = dataStore()
const state = ref({
  probeModal: null,
  newProbe: newProbe
})
const loadingText = "Loading..."

onMounted(() => {
    state.value.probeModal = new Modal('#probeModal', { keyboard: false, backdrop: "static" })
})

function showProbeModal() {
  state.value.newProbe = { ...newProbe }
  state.value.probeModal.show()
}

function editProbe(probe) {
  state.value.newProbe = { ...probe }
  state.value.probeModal.show()
}

async function saveProbe() {
  if (await dStore.saveProbes(state.value.newProbe)) {
    state.value.newProbe = newProbe
    state.value.probeModal.hide();
  }
  dStore.loadProbes()
}

// a computed ref
const loadedMessage = computed(() => {
  return dStore.fetchError ? dStore.fetchError.message : loadingText
})

</script>
<template>
<div>
  <p v-if="!dStore.policies">{{ loadedMessage }}</p>
  <div v-else>
    <button @click="showProbeModal()" class="btn btn-primary">Add Probe</button>
    <table class="table">
      <thead class="thead-dark">
        <tr>
          <th scope="col">#</th>
          <th scope="col">Collector</th>
          <th scope="col">Policy Name</th>
          <th scope="col">Policy Version</th>
          <th scope="col">Address</th>
          <th scope="col">Port</th>
          <th scope="col">User</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(probe, index) in dStore.probes" @click="editProbe(probe)">
          <th scope="row">{{index+1}}</th>
          <td>{{probe.collector}}</td>
          <td>{{dStore.policies[probe.policy].name}}</td>
          <td>{{probe.version}}</td>
          <td>{{probe.address}}</td>
          <td>{{probe.port}}</td>
          <td>{{probe.user}}</td>
        </tr>
      </tbody>
  </table>
</div>
  <div class="modal fade" id="probeModal" tabindex="-1" aria-labelledby="probeModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="probeModalLabel">Add probe</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveProbe()">
              <div class="mb-3">
                <label for="collectorInput" class="form-label">Collector</label>
                <select id="collectorInput" class="form-select" aria-label="Select collector" v-model="state.newProbe.collector">
                  <option v-for="coll,index in dStore.collectors" :value="index">{{index}}</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="policyInput" class="form-label">Policy Type</label>
                <select id="policyInput" class="form-select" aria-label="Select policy type" v-model="state.newProbe.policy">
                  <option v-for="(pol,key) in dStore.policies" :value="key">{{pol.name}}</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="versionInput" class="form-label">Version</label>
                <select id="versionInput" class="form-select" aria-label="Select policy version" v-model="state.newProbe.version">
                  <option v-if="state.newProbe.policy" v-for="version in dStore.policies[state.newProbe.policy].versions" :value="version">{{version}}</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="ipInput" class="form-label">Address</label>
                <input type="text" class="form-control" id="ipInput" aria-describedby="ipHelp" v-model="state.newProbe.address">
                <div id="ipHelp" class="form-text">IP or FQDN of the device</div>
              </div>
              <div class="mb-3">
                <label for="portInput" class="form-label">Port</label>
                <input type="text" class="form-control" id="portInput" aria-describedby="portHelp" v-model="state.newProbe.port">
                <div id="portHelp" class="form-text">Port where the service is listening</div>
              </div>
              <div class="mb-3">
                <label for="userInput" class="form-label">User</label>
                <input type="text" class="form-control" id="userInput" aria-describedby="userHelp" v-model="state.newProbe.user">
                <div id="userHelp" class="form-text">User with discover capabilities</div>
              </div>
              <div class="mb-3">
                <label for="passwordInput" class="form-label">Password</label>
                <input type="password" class="form-control" id="passwordInput" v-model="state.newProbe.password">
              </div>
              <button type="submit" class="btn btn-primary">Submit</button>
            </form>
          </div>
        </div>
      </div>
    </div>
</div>
</template>
