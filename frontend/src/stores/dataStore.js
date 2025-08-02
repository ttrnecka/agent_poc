import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { useApiStore } from '@/stores/apiStore'

export const useDataStore = defineStore('data', () => {
  const isLoggedIn = ref(true)
  const user = reactive({})
  
  const apiStore = useApiStore()

  async function getData() {
    await load("/user",user)
    if (isLoggedIn.value) {
      apiStore.reload()
    }
  }

  async function load(url, reactive, timeoutMs = 10000) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), timeoutMs);

    try {
      const res = await fetch(url, { signal: controller.signal });

      clearTimeout(timeout);

      if (!res.ok) {
        if (res.status == 401) {
          console.log("GOT 401")
          isLoggedIn.value = false;
        }
        throw new Error(`API service not available: HTTP status: ${res.status}`);
      }

      let data;
      try {
        data = await res.json();
      } catch (jsonError) {
        throw new Error("Failed to parse JSON response");
      }

      Object.assign(reactive, data)
      isLoggedIn.value = true;

    } catch (error) {
        console.error("Fetch failed:", error.name === 'AbortError' ? 'Request timed out' : error.message || error);
    }
  }
  return { isLoggedIn, getData }
});