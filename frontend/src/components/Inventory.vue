<template>
  <div class="container-fluid">
    <p v-if="!apiStore.collectors">{{ loadedMessage }}</p>
    <div v-else class="row">
      <div class="col-auto" style="flex: 0 0 200px;">
      </div>
      <div class="col">
        <div class="accordion" id="collectorsAccordion">
          <div
            class="accordion-item"
            v-for="(collector,collector_name) in apiStore.collectors"
            :key="collector_name"
          >
            <h2 class="accordion-header" :id="`heading-${collector_name}`">
              <button
                class="accordion-button collapsed"
                type="button"
                data-bs-toggle="collapse"
                :data-bs-target="`#collapse-${collector_name}`"
                aria-expanded="false"
                :aria-controls="`collapse-${collector_name}`"
              >
                {{ collector_name }}
              </button>
            </h2>
            <div
              :id="`collapse-${collector_name}`"
              class="accordion-collapse collapse"
              :aria-labelledby="`heading-${collector_name}`"
              data-bs-parent="#collectorsAccordion"
            >
              <div class="accordion-body">
                <Collector :collector="collector_name" :accordion-id="collector_name" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import Collector from './Collector.vue'
import { useApiStore } from '@/stores/apiStore'

const route = useRoute()
const apiStore = useApiStore()

const loadingText = "Loading..."

const loadedMessage = computed(() => {
  return apiStore.fetchError ? apiStore.fetchError.message : loadingText
})

// Watch route changes, reload only if we are on /inventory path
watch(
  () => route.fullPath,
  (newPath) => {
    if (newPath === '/inventory') {
      apiStore.loadCollectors()
    }
  }
)
</script>