<template>
  <div class="container-fluid">
    <p v-if="!apiStore.sortedCollectors">{{ loadedMessage }}</p>
    <div v-else class="row">
      <div class="col">
        <div class="accordion" id="collectorsAccordion">
          <div
            class="accordion-item"
            v-for="(collector,idx) in apiStore.collectors"
            :key="idx"
          >
            <h2 class="accordion-header" :id="`heading-${collector.name}`">
              <button
                class="accordion-button collapsed"
                type="button"
                data-bs-toggle="collapse"
                :data-bs-target="`#collapse-${collector.name}`"
                aria-expanded="false"
                :aria-controls="`collapse-${collector.name}`"
              >
                {{ collector.name }}
              </button>
            </h2>
            <div
              :id="`collapse-${collector.name}`"
              class="accordion-collapse collapse"
              :aria-labelledby="`heading-${collector.name}`"
              data-bs-parent="#collectorsAccordion"
            >
              <div class="accordion-body">
                <Collector :collector="collector.name" :accordion-id="collector.name" />
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