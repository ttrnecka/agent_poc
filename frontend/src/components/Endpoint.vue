<template>
  <div>
    <div v-if="loading">Loading endpoint data...</div>
    <div v-else-if="error" class="text-danger">{{ error }}</div>
    <pre v-else class="bg-light p-3 rounded" style="white-space: pre-wrap;">{{ content }}</pre>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApiStore } from '@/stores/apiStore'

const apiStore = useApiStore()

const props = defineProps({
  collector: String,
  device: String,
  endpoint: String,
})

const content = ref('')
const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  try {
    const res = await fetch(apiStore.endpointEndpoint(props.collector,props.device,props.endpoint))
    if (!res.ok) throw new Error('Failed to load endpoint data')
    content.value = await res.text()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>
