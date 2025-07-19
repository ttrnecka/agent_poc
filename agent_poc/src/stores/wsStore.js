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
      console.log('📥 Received data:', evt.data);

      const messages = evt.data
      .split('\n')
      .filter(line => line.trim().length > 0);

      for (const msgStr of messages) {
        try {
          console.log('📥 Received json:', msgStr);
          const message = JSON.parse(msgStr);
          handleMessage(message);
        } catch (e) {
          console.error("❌ Failed to parse message:", msgStr, e);
        }
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

  function handleMessage(msg) {
    console.log('📥 Decoded JSON message:', msg);
    const clientId = getClientId();
    if (msg.Type == 1 && msg.Source != clientId) {
      apiStore.updateCollectorState(msg.Source,{status: msg.Text})
    }
  }
  
  return { conn, wsUrl, connect }
});