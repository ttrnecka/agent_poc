import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useApiStore = defineStore('api', () => {
  const policies = ref(null)
  const probes = ref(null)
  const collectors = ref(null)
  const fetchError = ref(null)

  const httpUrl = ref(import.meta.env.VITE_HTTP_PROTOCOL + '://' + import.meta.env.VITE_API_HOST);

  const policiesUrl = computed(() => {
    return `${httpUrl.value}${import.meta.env.VITE_APP_POLICY_ENDPOINT}`; 
  })

  const probesUrl = computed(() => {
    return `${httpUrl.value}${import.meta.env.VITE_APP_PROBE_ENDPOINT}`; 
  })

  const collectorsUrl = computed(() => {
    return `${httpUrl.value}${import.meta.env.VITE_APP_COLLECTOR_ENDPOINT}`;
  })
  
  async function load(url,ref) {
    try {  
      const res = await fetch(
        url
      )
      ref.value = await res.json()
    } catch (error) {
      console.error("Error:", error.message);
      fetchError.value = error;
    }
  }
  async function loadPolicies() {
    await load(policiesUrl.value,policies)
  }
  
  async function loadProbes() {
    await load(probesUrl.value,probes)
  }

  async function loadCollectors() {
    await load(collectorsUrl.value,collectors)
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
      const response = await fetch(probesUrl(), {
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
    collectors.value[collector] = state
  }


  return { policies, loadPolicies, probes, loadProbes, saveProbes, fetchError, collectors, loadCollectors, updateCollectorState, httpUrl, policiesUrl, probesUrl, collectorsUrl }
});