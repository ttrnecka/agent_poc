<template>
  <div class="modal show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5)">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Login</h5>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitLogin">
            <div class="mb-3">
              <label>Login</label>
              <input v-model="username" type="text" class="form-control" required />
            </div>
            <div class="mb-3">
              <label>Password</label>
              <input v-model="password" type="password" class="form-control" required />
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useDataStore } from '@/stores/dataStore'

const dataStore = useDataStore()


const username = ref('')
const password = ref('')
const router = useRouter()

const submitLogin = async () => {
  try {
    const formData = new FormData()
    formData.append('username', username.value)
    formData.append('password', password.value)
    const res = await fetch('/login', {
      method: 'POST',
      body: formData,
    })

    if (res.ok) {
      const data = await res.json()
      dataStore.setLoggedIn(true)
      router.push('/')
    } else {
      alert('Invalid login')
    }
   } catch (err) {
    console.error(err)
    alert('Error logging in')
  }
}
</script>
