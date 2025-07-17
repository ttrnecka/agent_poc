import { defineStore } from 'pinia'
import { ref } from 'vue'

export const dataStore = defineStore('defStore', () => {
  const policies = ref(null)
  const probes = ref(null)
  const collectors = ref(null)
  const fetchError = ref(null)
  
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
    await load(`http://localhost:8888/api/v1/policy`,policies)
  }
  
  async function loadProbes() {
    await load(`http://localhost:8888/api/v1/probe`,probes)
  }

  async function loadCollectors() {
    await load(`http://localhost:8888/api/v1/collector`,collectors)
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
      const response = await fetch("http://localhost:8888/api/v1/probe", {
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

  return { policies, loadPolicies, probes, loadProbes, saveProbes, fetchError, collectors, loadCollectors, updateCollectorState }
})