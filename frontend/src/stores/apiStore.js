import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import axios from 'axios'

const POLICY_ENDPOINT="/api/v1/policy"
const PROBE_ENDPOINT="/api/v1/probe"
const COLLECTORS_ENDPOINT="/api/v1/collector"
const COLLECTOR_ENDPOINT="/api/v1/data/collector/"

const collectorEndpoint = (collector) => `${COLLECTOR_ENDPOINT}${collector}`
const deviceEndpoint = (collector,device) => `${COLLECTOR_ENDPOINT}${collector}/${device}`
const endpointEndpoint = (collector,device,endpoint) => `${COLLECTOR_ENDPOINT}${collector}/${device}/${endpoint}`

export const useApiStore = defineStore('api', () => {
  const policies = ref([])
  const probes = ref([])
  const collectors = ref([])
  const fetchError = ref(null)
  
  const sortedCollectors = computed(() => collectors.value.sort((a, b) => a.name.localeCompare(b.name, undefined, { sensitivity: 'base' })))

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

  async function deleteProbe(probeId) {
    try {
      await axios.delete(`${PROBE_ENDPOINT}/${probeId}`)
      probes.value = probes.value.filter(obj => obj.id !== probeId)
    }
    catch (error) {
      console.error("Error:", error);
      return false
    }
  }

  async function deleteCollector(collId) {
    try {
      await axios.delete(`${COLLECTORS_ENDPOINT}/${collId}`)
      // collectors.value = collectors.value.filter(obj => obj.id !== collId)
      reload()
    }
    catch (error) {
      console.error("Error:", error);
      return false
    }
  }

  async function saveCollector(collector) {
    if (collector.id) {
      collector = await post(`${COLLECTORS_ENDPOINT}/${collector.id}`, collector)
      if (collector) {
        // for (let i in collectors.value) {
        //   if (collectors.value[i].id == collector.id) {
        //     collectors.value[i] = collector
        //   }
        // }
        reload()
        return true
      }
      return false
    }
    // new probe
    collector = await post(COLLECTORS_ENDPOINT, collector)
    if (collector) { 
      // collectors.value.push(collector)
      reload()
      return true
    }
    return false
  }

  async function saveProbe(probe) {
    if (probe.id) {
      probe = await post(`${PROBE_ENDPOINT}/${probe.id}`, probe)
      if (probe) {
        // for (let i in probes.value) {
        //   if (probes.value[i].id == probe.id) {
        //     probes.value[i] = probe
        //   }
        // }
        reload()
        return true
      }
      return false
    }
    // new probe
    probe = await post(PROBE_ENDPOINT, probe)
    if (probe) { 
      // probes.value.push(probe)
      reload()
      return true
    }
    return false
  }

  async function post(endpoint,data) {
    try {
      const response = await axios.post(endpoint, data, {
        headers: {
          "Content-Type": "application/json",
        },
      });
  
      console.log("Success:", response.status);
      return response.data
    } catch (error) {
      console.error("Error:", error);
      return false
    }
  }

  function updateCollectorStatus(collector,status) {
    collectors.value && (collectors.value.find((elm) => elm.name === collector).status = status);
  }


  function getCollector(id) {
    return collectors.value.find((o) => o.id === id)
  }
  return { policies, loadPolicies, probes, loadProbes, saveProbe, saveCollector, fetchError, collectors, sortedCollectors, loadCollectors, updateCollectorStatus,
           endpointEndpoint, deviceEndpoint, collectorEndpoint, reload, getCollector, deleteProbe,deleteCollector }
});
