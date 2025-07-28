<template>
  <div class="container mt-4">
    <div class="accordion" id="collectorsAccordion">
      <div
        class="accordion-item"
        v-for="collector in collectors"
        :key="collector"
      >
        <h2 class="accordion-header" :id="`heading-${collector}`">
          <button
            class="accordion-button collapsed"
            type="button"
            data-bs-toggle="collapse"
            :data-bs-target="`#collapse-${collector}`"
            aria-expanded="false"
            :aria-controls="`collapse-${collector}`"
          >
            {{ collector }}
          </button>
        </h2>
        <div
          :id="`collapse-${collector}`"
          class="accordion-collapse collapse"
          :aria-labelledby="`heading-${collector}`"
          data-bs-parent="#collectorsAccordion"
        >
          <div class="accordion-body">
            <Collector :collector="collector" :accordion-id="collector" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import Collector from './Collector.vue'

const collectors = ref([])
const route = useRoute()

const fetchCollectors = async () => {
  try {
    const res = await fetch('/api/v1/data/collector')
    if (!res.ok) throw new Error('Failed to load collectors')
    collectors.value = await res.json()
  } catch (err) {
    console.error(err)
  }
}

// Initial fetch on mount
onMounted(fetchCollectors)

// Watch route changes, reload only if we are on /inventory path
watch(
  () => route.fullPath,
  (newPath) => {
    if (newPath === '/inventory') {
      fetchCollectors()
    }
  }
)
</script>