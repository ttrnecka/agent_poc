<script setup>
import { ref, onMounted, computed } from 'vue'
import { useApiStore } from '@/stores/apiStore'
import { Modal } from "bootstrap";

const newCollector = {
    id: null,
    name: null,
    status: null,
    password: null
}

const collectorDefaultStatus = "UNKNOWN"
const loadingText = "Loading..."
const apiStore = useApiStore()

const state = ref({
  collectorModal: null,
  newCollector: newCollector
})

const loadedMessage = computed(() => {
  return apiStore.fetchError ? apiStore.fetchError.message : loadingText
})

function showCollectorModal() {
  state.value.newCollector = { ...newCollector }
  state.value.collectorModal.show()
}

onMounted(() => {
    state.value.collectorModal = new Modal('#collectorModal', { keyboard: false, backdrop: "static" })
})

function editCollector(collector) {
  state.value.newCollector = { ...collector }
  state.value.collectorModal.show()
}

async function saveCollector() {
  if (await apiStore.saveCollector(state.value.newCollector)) {
    state.value.newCollector = newCollector
    state.value.collectorModal.hide();
  }
  apiStore.loadProbes()
}
</script>
<template>
<div class="container-fluid">
  <p v-if="!apiStore.sortedCollectors">{{ loadedMessage }}</p>
  <div v-else class="row">
    <div class="col-auto" style="flex: 0 0 200px;">
      <button @click="showCollectorModal()" class="btn btn-primary btn-sm w-100">Add Collector</button>
    </div>
    <div class="col">
      <table class="table">
        <thead class="thead-dark">
          <tr>
            <th scope="col">Collector Name</th>
            <th scope="col">Status</th>
            <th scope="col">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(coll, index) in apiStore.collectors" @click="editCollector(coll)" :key="index" class="coll-row">
            <td>{{coll.name}}</td>
            <td>{{coll.status || collectorDefaultStatus}}</td>
            <td>
              <div class="d-flex gap-2">
                <button
                  @click.stop="apiStore.deleteCollector(coll.id)"
                  class="btn btn-primary btn-sm"
                >
                  Delete
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  <div class="modal fade" id="collectorModal" tabindex="-1" aria-labelledby="probeModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h1 class="modal-title fs-6" id="probeModalLabel">Add collector</h1>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveCollector()">
            <div class="mb-3">
              <input type="text" class="form-control form-control-sm" id="nameInput" aria-describedby="nameHelp" v-model="state.newCollector.name"
              placeholder="Collector name" title="Collector name">
            </div>
            <div class="mb-3">
              <input type="password" class="form-control form-control-sm" id="passwordInput" v-model="state.newCollector.password" placeholder="Password"
              title="Password">
            </div>
            <button type="submit" class="btn btn-primary btn-sm">Submit</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<style type="text/css">
.coll-row {
  cursor: pointer;
  transition: background 0.2s;
}
.coll-row:hover td {
  background: #cecece;
}
</style>