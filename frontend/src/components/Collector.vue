<template>
  <div class="accordion" :id="`devicesAccordion-${collector}`">
    <div
      class="accordion-item"
      v-for="device in devices"
      :key="device"
    >
      <h2 class="accordion-header" :id="`heading-${collector}-${device}`">
        <button
          class="accordion-button collapsed"
          type="button"
          data-bs-toggle="collapse"
          :data-bs-target="`#collapse-${collector}-${device}`"
          aria-expanded="false"
          :aria-controls="`collapse-${collector}-${device}`"
        >
          {{ device }}
        </button>
      </h2>
      <div
        :id="`collapse-${collector}-${device}`"
        class="accordion-collapse collapse"
        :aria-labelledby="`heading-${collector}-${device}`"
        :data-bs-parent="`#devicesAccordion-${collector}`"
      >
        <div class="accordion-body">
          <Device :collector="collector" :device="device" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Device from './Device.vue'
import { useApiStore } from '@/stores/apiStore'

const apiStore = useApiStore()

const props = defineProps({
  collector: String,
  accordionId: String
})

const devices = ref([])

onMounted(async () => {
  try {
    const res = await fetch(apiStore.collectorEndpoint(props.collector))
    if (!res.ok) throw new Error(`Failed to load devices for ${props.collector}`)
    const data = await res.json()
    devices.value = data["devices"] || []
  } catch (err) {
    console.error(err)
  }
})

</script>
