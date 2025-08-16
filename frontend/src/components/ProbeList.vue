<script setup>
import { ref, onMounted, computed, watch  } from 'vue'
import { useApiStore } from '@/stores/apiStore'
import { sendMessage, MESSAGE_TYPE } from '@/services/messages'
import { Modal } from "bootstrap";
import { useSessionStore } from '@/stores/sessionStore'

const newProbe = {
    id: null,
    policy: "",
    collector_id: "",
    version: "",
    address: null,
    port: null,
    user: null,
    password: null
}
const newProbeState = {
  errors: {},
  touched: {}
}
const apiStore = useApiStore()
const sessionStore = useSessionStore()

const state = ref({
  probeModal: null,
  taskModal: null,
  newProbe: structuredClone(newProbe),
  newProbeState: structuredClone(newProbeState),
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
  state.value.newProbe = structuredClone(newProbe)
  state.value.newProbeState = structuredClone(newProbeState)
  state.value.probeModal.show()
}

function editProbe(probe) {
  state.value.newProbe = { ...probe }
  state.value.newProbeState = structuredClone(newProbeState)
  state.value.probeModal.show()
}

async function saveProbe() {
  if (!validateProbeForm()) return

  if (await apiStore.saveProbe(state.value.newProbe)) {
    state.value.newProbe = newProbe
    state.value.probeModal.hide();
  }
  apiStore.loadProbes()
}

function runProbe(probe) {
  const session = sendMessage(MESSAGE_TYPE.PROBE_START, apiStore.getCollector(probe.collector_id).name, `PROBE_ID="${probe.id}" CLI_USER="${probe.user}" CLI_PASSWORD="${probe.password}" ${probe.policy}_${probe.version} collect --endpoint ${probe.address}:${probe.port}`);
  state.value.task.title = `${apiStore.getCollector(probe.collector_id).name} - ${probe.policy}_${probe.version} - ${probe.address}:${probe.port} - collection`;
  state.value.task.message = "Waiting for data...";
  state.value.task.state = "RUNNING";
  state.value.task.session = session;
  state.value.taskModal.show();
}

function validateProbe(probe) {
  const session = sendMessage(MESSAGE_TYPE.PROBE_START, apiStore.getCollector(probe.collector_id).name, `PROBE_ID="${probe.id}" CLI_USER="${probe.user}" CLI_PASSWORD="${probe.password}" ${probe.policy}_${probe.version} validate --endpoint ${probe.address}:${probe.port}`);
  state.value.task.title = `${apiStore.getCollector(probe.collector_id).name} - ${probe.policy}_${probe.version} - ${probe.address}:${probe.port} - validation`;
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
    if (newData.Type == MESSAGE_TYPE.PROBE_FINISHED_ERR) {
      state.value.task.state = "ERROR";
    }
    if (newData.Type == MESSAGE_TYPE.PROBE_FINISHED_OK) {
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

function validateProbeForm() {
  validateProbeField("collector_id")
  validateProbeField("policy")
  validateProbeField("version")
  validateProbeField("address")
  validateProbeField("port")
  validateProbeField("user")
  return  !isValid('collector_id') && 
          !isValid('policy') && 
          !isValid('version') && 
          !isValid('address') && 
          !isValid('port') && 
          !isValid('user')
}

function required(field,name="Field") {
  const errors = state.value.newProbeState.errors
  const form = state.value.newProbe
  if (!form[field]) {
    errors[field] = `${name} is required`
  } else {
    errors[field] = ""
  }
}
function validateProbeField(field) {
  state.value.newProbeState.touched[field] = true
  const errors = state.value.newProbeState.errors
  const form = state.value.newProbe
  required(field)
}

function isValid(field) {
  return state.value.newProbeState.errors[field] ? true : false
}

function isInvalid(field) {
  return state.value.newProbeState.touched[field] && !state.value.newProbeState.errors[field]
}

function invalidError(field) {
  return state.value.newProbeState.errors[field]
}

</script>
<template>
<div class="container-fluid">
  <p v-if="!apiStore.policies">{{ loadedMessage }}</p>
  <div v-else class="row">
    <div class="col-auto" style="flex: 0 0 200px;">
      <button @click="showProbeModal()" class="btn btn-primary btn-sm w-100">Add Probe</button>
    </div>
    <div class="col">
      <table class="table table-striped table-hover table-sm">
        <thead class="">
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
            <td>{{apiStore.getCollector(probe.collector_id).name}}</td>
            <td>{{apiStore.policies.find((elm) => elm.name === probe.policy).description}}</td>
            <td>{{probe.version}}</td>
            <td>{{probe.address}}</td>
            <td>{{probe.port}}</td>
            <td>{{probe.user}}</td>
            <td>
              <div class="d-flex gap-2">
                <button
                  @click.stop="runProbe(probe)"
                  class="btn btn-primary btn-sm"
                >
                  Run
                </button>
                <button
                  @click.stop="validateProbe(probe)"
                  class="btn btn-primary btn-sm"
                >
                  Validate
                </button>
                <button
                  @click.stop="apiStore.deleteProbe(probe.id)"
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
  <div>
    
</div>
  <div class="modal fade" id="probeModal" tabindex="-1" aria-labelledby="probeModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-6" id="probeModalLabel">Add probe</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveProbe()">
              <div class="mb-3">
                <select   id="collectorInput" 
                          class="form-select form-select-sm" 
                          aria-label="Select collector" 
                          v-model="state.newProbe.collector_id"
                          :class="{'is-invalid': isValid('collector_id'), 'is-valid': isInvalid('collector_id')}"
                          @blur="validateProbeField('collector_id')"
                          >
                  <option selected disabled value="">-- Collector --</option>
                  <option v-for="coll,index in apiStore.sortedCollectors" :value="coll.id" :key="index">{{coll.name}}</option>
                </select>
                <div v-if="isValid('collector_id')" class="invalid-feedback">{{ invalidError('collector_id') }}</div>
              </div>
              <div class="mb-3">
                <select id="policyInput" 
                        class="form-select form-select-sm" 
                        aria-label="Select policy type" 
                        v-model="state.newProbe.policy"
                        :class="{'is-invalid': isValid('policy'), 'is-valid': isInvalid('policy')}"
                        @blur="validateProbeField('policy')">
                  <option selected disabled value="">-- Policy --</option>
                  <option v-for="(pol,key) in apiStore.policies" :value="pol.name" :key="key">{{pol.description}}</option>
                </select>
                <div v-if="isValid('policy')" class="invalid-feedback">{{ invalidError('policy') }}</div>
              </div>
              <div class="mb-3">
                <select
                  id="versionInput"
                  class="form-select form-select-sm"
                  aria-label="Select policy version"
                  v-model="state.newProbe.version"
                  :class="{'is-invalid': isValid('version'), 'is-valid': isInvalid('version')}"
                  @blur="validateProbeField('version')"
                >
                  <option v-if="!state.newProbe.policy" disabled value="">-- Select policy first --</option>
                  <option v-else disabled value="">-- Version --</option>
                  <option
                  v-for="version in (state.newProbe.policy ? apiStore.policies.find((elm) => elm.name === state.newProbe.policy).versions : [])"
                  :value="version"
                  :key="version"
                  >
                  {{ version }}
                  </option>
                </select>
                <div v-if="isValid('version')" class="invalid-feedback">{{ invalidError('version') }}</div>
              </div>
              <div class="mb-3">
                <input  type="text" 
                        class="form-control form-control-sm" 
                        id="ipInput" 
                        aria-describedby="ipHelp" 
                        v-model="state.newProbe.address"
                        placeholder="IP or FQDN of the device"
                        :class="{'is-invalid': isValid('address'), 'is-valid': isInvalid('address')}"
                        @blur="validateProbeField('address')"
                        title="IP or FQDN of the device">
                <div v-if="isValid('address')" class="invalid-feedback">{{ invalidError('address') }}</div>
              </div>
              <div class="mb-3">
                <input  type="number" 
                        class="form-control form-control-sm" 
                        id="portInput" 
                        aria-describedby="portHelp" 
                        v-model="state.newProbe.port"
                        placeholder="Port" 
                        :class="{'is-invalid': isValid('port'), 'is-valid': isInvalid('port')}"
                        @blur="validateProbeField('port')"
                        title="Port">
                <div v-if="isValid('port')" class="invalid-feedback">{{ invalidError('port') }}</div>
              </div>
              <div class="mb-3">
                <input  type="text"  
                        class="form-control form-control-sm" 
                        id="userInput" 
                        aria-describedby="userHelp" 
                        v-model="state.newProbe.user"
                        placeholder="User with discover capabilities" 
                        :class="{'is-invalid': isValid('user'), 'is-valid': isInvalid('user')}"
                        @blur="validateProbeField('user')"
                        title="User with discover capabilities" >
                <div v-if="isValid('user')" class="invalid-feedback">{{ invalidError('user') }}</div>
              </div>
              <div class="mb-3">
                <input type="password" class="form-control form-control-sm" id="passwordInput" v-model="state.newProbe.password" placeholder="Password"
                title="Password">
              </div>
              <button type="submit" class="btn btn-primary btn-sm">Submit</button>
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