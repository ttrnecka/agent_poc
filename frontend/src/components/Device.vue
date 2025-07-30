<template>
  <div class="accordion" :id="`endpointsAccordion-${collector}-${device}`">
    <div
      class="accordion-item"
      v-for="endpoint in endpoints"
      :key="endpoint"
    >
      <h2 class="accordion-header" :id="`heading-${collector}-${device}-${endpoint}`">
        <button
          class="accordion-button collapsed"
          type="button"
          data-bs-toggle="collapse"
          :data-bs-target="`#collapse-${collector}-${device}-${endpoint}`"
          aria-expanded="false"
          :aria-controls="`collapse-${collector}-${device}-${endpoint}`"
        >
          {{ endpoint }}
        </button>
      </h2>
      <div
        :id="`collapse-${collector}-${device}-${endpoint}`"
        class="accordion-collapse collapse"
        :aria-labelledby="`heading-${collector}-${device}-${endpoint}`"
        :data-bs-parent="`#endpointsAccordion-${collector}-${device}`"
      >
        <div class="accordion-body">
          <Endpoint :collector="collector" :device="device" :endpoint="endpoint" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Endpoint from './Endpoint.vue'
import { useApiStore } from '@/stores/apiStore'

const apiStore = useApiStore()

const props = defineProps({
  collector: String,
  device: String,
})

const endpoints = ref([])

onMounted(async () => {
  try {
        const res = await fetch(apiStore.deviceEndpoint(props.collector,props.device))
    if (!res.ok) throw new Error(`Failed to load endpoints for ${props.device}`)
    const data = await res.json()
    endpoints.value = data["endpoints"] || []
  } catch (err) {
    console.error(err)
  }
})
 
</script>
