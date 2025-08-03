import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

const POLICY_ENDPOINT="/api/v1/policy"
const PROBE_ENDPOINT="/api/v1/probe"
const COLLECTORS_ENDPOINT="/api/v1/collector"
const COLLECTOR_ENDPOINT="/api/v1/data/collector/"

const collectorEndpoint = (collector) => `${COLLECTOR_ENDPOINT}${collector}`
const deviceEndpoint = (collector,device) => `${COLLECTOR_ENDPOINT}${collector}/${device}`
const endpointEndpoint = (collector,device,endpoint) => `${COLLECTOR_ENDPOINT}${collector}/${device}/${endpoint}`

export const useApiStore = defineStore('api', () => {
  const policies = ref(null)
  const probes = ref(null)
  const collectors = ref(null)
  const fetchError = ref(null)
  
  function reload() {
    loadCollectors()
    loadPolicies()
    loadProbes()
  }

  async function load(url, ref, timeoutMs = 10000) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), timeoutMs);

    try {
      const res = await fetch(url, { signal: controller.signal });

      clearTimeout(timeout);

      if (!res.ok) {
        throw new Error(`API service not available: HTTP status: ${res.status}`);
      }

      let data;
      try {
        data = await res.json();
      } catch (jsonError) {
        throw new Error("Failed to parse JSON response");
      }

      ref.value = data;

    } catch (error) {
        console.error("Fetch failed:", error.name === 'AbortError' ? 'Request timed out' : error.message || error);
        fetchError.value = error;
    }
  }

  async function loadPolicies() {
    await load(POLICY_ENDPOINT,policies)
  }
  
  async function loadProbes() {
    await load(PROBE_ENDPOINT,probes)
  }

  async function loadCollectors() {
    await load(COLLECTORS_ENDPOINT,collectors)
  }

  async function saveProbes(data) {
    try {
      if (!probes.value) {
        probes.value = []
      }
      // new probe
      if (!data.id) {
        probes.value.push(data)
      } // update probe
      else {
        for (let i in probes.value) {
          if (probes.value[i].id == data.id) {
            probes.value[i] = data
          }
        }
      }
      const response = await fetch(PROBE_ENDPOINT, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(probes.value),
      });
  
      if (!response.ok) {
        throw new Error(`Error! status: ${response.status}`);
      }
      console.log("Success:", response.status);
    } catch (error) {
      console.error("Error:", error);
      return false
    }
    return true
  }

  function updateCollectorState(collector,state) {
    collectors.value && (collectors.value[collector] = state);
  }


  return { policies, loadPolicies, probes, loadProbes, saveProbes, fetchError, collectors, loadCollectors, updateCollectorState,
           endpointEndpoint, deviceEndpoint, collectorEndpoint, reload }
});