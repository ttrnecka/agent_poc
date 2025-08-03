import { defineStore } from 'pinia'
import { ref, reactive,computed } from 'vue'
import { useApiStore } from '@/stores/apiStore'
import { useRouter } from 'vue-router'

export const useDataStore = defineStore('data', () => {
  const loggedIn = ref(false)
  const user = reactive({})
  
  const apiStore = useApiStore()

  const isLoggedIn = computed(() => loggedIn.value)

  const router = useRouter()

  function setLoggedIn(value) {
    loggedIn.value = value
  }

  async function getData() {
    // await load("/user",user)
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
          loggedIn.value = false;
          router.push("/login")
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
      loggedIn.value = true;

    } catch (error) {
        console.error("Fetch failed:", error.name === 'AbortError' ? 'Request timed out' : error.message || error);
    }
  }

  //initialize the user and set the loggedIn var
  load("/user",user)

  return { loggedIn,user, isLoggedIn, getData, setLoggedIn}
});