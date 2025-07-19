import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useApiStore } from '@/stores/apiStore'

export const useWsConnectionStore = defineStore('wsConnection', () => {
  const conn = ref(null) 
  const wsUrl = ref(import.meta.env.VITE_WS_PROTOCOL + '://' + import.meta.env.VITE_API_HOST + '/ws');

  const apiStore = useApiStore();
  
  function connect() {
    conn.value = new WebSocket(wsUrl.value);

    conn.value.onmessage = function (evt) {
      console.log(evt.data)
      const clientId = getClientId();
      const message = JSON.parse(evt.data)
      if (message['Type'] == 1 && message['Source'] != clientId) {
        apiStore.updateCollectorState(message['Source'],{status: message['Text']})
      }
    };

    conn.value.onopen = function (evt) {
      const clientId = getClientId();
      const msg = {
        Type: 1,
        Source: clientId,
        Text: "web_client connected"
      };
      conn.value.send(JSON.stringify(msg));
    };
  }

  function getClientId() {
    let clientId = localStorage.getItem("client_id");
    if (!clientId) {
      clientId = "web_client_" + crypto.randomUUID(); // Secure, random
      localStorage.setItem("client_id", clientId);
    }
    return clientId;
  }
  return { conn, wsUrl, connect }
});