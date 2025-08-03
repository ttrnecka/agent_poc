<template>
  <div class="modal show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title fs-6">Login</h5>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitLogin">
            <div class="mb-3">
              <input ref="focusInput" v-model="username" type="text" class="form-control form-control-sm" placeholder="Login" required />
            </div>
            <div class="mb-3">
              <input v-model="password" type="password" class="form-control form-control-sm" placeholder="Password" required />
            </div>
            <div class="modal-footer">
              <button type="submit" class="btn btn-primary">Login</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useDataStore } from '@/stores/dataStore'

const dataStore = useDataStore()

const username = ref('')
const password = ref('')
const focusInput = ref(null)

const submitLogin = async () => {
  await dataStore.loginUser(username.value,password.value)
}

onMounted(() => {
  focusInput.value?.focus()
})
</script>
