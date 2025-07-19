import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useApiStore } from '@/stores/apiStore'

const MESSAGE_TYPE = {
  'ONLINE': 1,
  'OFFLINE': 2,
  'REFRESH': 10
};

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
    console.log('📥 Message:', msg.Text);
    if (msg.Source != getClientId()) {
      if (msg.Type == MESSAGE_TYPE.ONLINE) {
        apiStore.updateCollectorState(msg.Source,{status: "ONLINE"})
      }
      if (msg.Type == MESSAGE_TYPE.OFFLINE) {
        apiStore.updateCollectorState(msg.Source,{status: "OFFLINE"})
      }
    }
  }
  
  return { conn, wsUrl, connect }
});