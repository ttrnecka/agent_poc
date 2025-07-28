<script setup>
import { ref, onMounted, computed, watch  } from 'vue'
import { useApiStore } from '@/stores/apiStore'
import { useWsConnectionStore } from '@/stores/wsStore'
import { Modal } from "bootstrap";
import { MESSAGE_TYPE } from '@/stores/messages'
import { useSessionStore } from '@/stores/sessionStore'

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
const apiStore = useApiStore()
const ws = useWsConnectionStore()
const sessionStore = useSessionStore()

const state = ref({
  probeModal: null,
  taskModal: null,
  newProbe: newProbe,
  task: {
    title: null,
    message: null,
    session: null,
    state: "WAITING" // "WAITING", "RUNNING", "FINISHED", "ERROR"
  }
})
const loadingText = "Loading..."

onMounted(() => {
    state.value.probeModal = new Modal('#probeModal', { keyboard: false, backdrop: "static" })
    state.value.taskModal = new Modal('#taskModal', { keyboard: false, backdrop: "static" })
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
  if (await apiStore.saveProbes(state.value.newProbe)) {
    state.value.newProbe = newProbe
    state.value.probeModal.hide();
  }
  apiStore.loadProbes()
}

function runProbe(probe) {
  const session = ws.sendMessage(MESSAGE_TYPE.RUN, probe.collector, `PROBE_ID="${probe.id}" CLI_USER="${probe.user}" CLI_PASSWORD="${probe.password}" ${probe.policy}_${probe.version} collect --endpoint ${probe.address}:${probe.port}`);
  state.value.task.title = `${probe.collector} - ${probe.policy}_${probe.version} - ${probe.address}:${probe.port} - collection`;
  state.value.task.message = "Waiting for data...";
  state.value.task.state = "RUNNING";
  state.value.task.session = session;
  state.value.taskModal.show();
}

function validateProbe(probe) {
  const session = ws.sendMessage(MESSAGE_TYPE.RUN, probe.collector, `PROBE_ID="${probe.id}" CLI_USER="${probe.user}" CLI_PASSWORD="${probe.password}" ${probe.policy}_${probe.version} validate --endpoint ${probe.address}:${probe.port}`);
  state.value.task.title = `${probe.collector} - ${probe.policy}_${probe.version} - ${probe.address}:${probe.port} - validation`;
  state.value.task.message = "Waiting for data...";
  state.value.task.state = "RUNNING";
  state.value.task.session = session;
  state.value.taskModal.show();
}
// a computed ref
const loadedMessage = computed(() => {
  return apiStore.fetchError ? apiStore.fetchError.message : loadingText
})

const sessionData = computed(() => sessionStore.sessions[state.value.task.session]);

watch(sessionData, (newData) => {
  if (newData) {
    console.log(`New message for session ${state.value.task.session}:`, newData);
    state.value.task.message = newData.Text;
    if (newData.Type == MESSAGE_TYPE.FINISHED_ERR) {
      state.value.task.state = "ERROR";
    }
    if (newData.Type == MESSAGE_TYPE.FINISHED_OK) {
      state.value.task.state = "FINISHED";
    }
    // Handle message logic here
    
    // attach this to onclose on tha taskModal
    // sessionStore.clearSession(state.value.task.session); // optional
  }
});

const stateClass = computed(() => {
  switch (state.value.task.state) {
    case "FINISHED":
      return "text-success";
    case "ERROR":
      return "text-danger";
    default:
      return "";
  }
});

</script>
<template>
<div>
  <p v-if="!apiStore.policies">{{ loadedMessage }}</p>
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
          <th scope="col">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(probe, index) in apiStore.probes" @click="editProbe(probe)" :key="index" class="probe-row">
          <th scope="row">{{index+1}}</th>
          <td>{{probe.collector}}</td>
          <td>{{apiStore.policies[probe.policy].name}}</td>
          <td>{{probe.version}}</td>
          <td>{{probe.address}}</td>
          <td>{{probe.port}}</td>
          <td>{{probe.user}}</td>
          <td>
            <div class="d-flex gap-2">
              <button
                @click.stop="runProbe(probe)"
                class="btn btn-primary"
              >
                Run
              </button>
              <button
                @click.stop="validateProbe(probe)"
                class="btn btn-primary"
              >
                Validate
              </button>
            </div>
          </td>
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
                  <option v-for="coll,index in apiStore.collectors" :value="index" :key="index">{{index}}</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="policyInput" class="form-label">Policy Type</label>
                <select id="policyInput" class="form-select" aria-label="Select policy type" v-model="state.newProbe.policy">
                  <option v-for="(pol,key) in apiStore.policies" :value="key" :key="key">{{pol.name}}</option>
                </select>
              </div>
                <div class="mb-3">
                <label for="versionInput" class="form-label">Version</label>
                <select
                  id="versionInput"
                  class="form-select"
                  aria-label="Select policy version"
                  v-model="state.newProbe.version"
                >
                  <option v-if="!state.newProbe.policy" disabled value="">-- Select policy first --</option>
                  <option
                  v-for="version in (state.newProbe.policy ? apiStore.policies[state.newProbe.policy].versions : [])"
                  :value="version"
                  :key="version"
                  >
                  {{ version }}
                  </option>
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

    <div class="modal fade" id="taskModal" tabindex="-1" aria-labelledby="taskModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header" :class="stateClass">
            <h1 class="modal-title fs-5" id="taskModalLabel">Task: {{ state.task.title }}</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <pre>{{ state.task.message }}</pre>
          </div>
        </div>
      </div>
    </div>
</div>
</template>

<style type="text/css">
.probe-row {
  cursor: pointer;
  transition: background 0.2s;
}
.probe-row:hover td {
  background: #cecece;
}
.text-success {
  background-color: rgb(138, 237, 138);
}
.text-danger {
  background-color: rgb(255 194 194);
}
</style>