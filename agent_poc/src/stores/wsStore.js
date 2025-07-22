import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useApiStore } from '@/stores/apiStore'
import { MESSAGE_TYPE } from '@/stores/messages'
import { useSessionStore } from '@/stores/sessionStore';


export const useWsConnectionStore = defineStore('wsConnection', () => {
  const conn = ref(null) 
  const wsUrl = ref(import.meta.env.VITE_WS_PROTOCOL + '://' + import.meta.env.VITE_API_HOST + '/ws');

  const apiStore = useApiStore();
  const store = useSessionStore();
  
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

    
    conn.value.addEventListener('message', (event) => {
      try {
        const data = JSON.parse(event.data);
        store.handleWebSocketMessage(data);
      } catch (err) {
        console.error('Invalid WebSocket message:', err);
      }
    });

    conn.value.onopen = function (evt) {
      sendMessage(MESSAGE_TYPE.ONLINE, null, "web_client connected", null);
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

  /**
   * Sends a message over the WebSocket connection.
   * @param {string} type - The message type.
   * @param {string|null} destinationId - The destination client ID, or null if not applicable.
   * @param {string} text - The message text.
   * @param {string} [session=crypto.randomUUID()] - The session ID.
   * @returns {string} The session ID used for the message.
   */
  function sendMessage(type, destinationId, text, session = crypto.randomUUID()) {
    const clientId = getClientId(); 
    const msg = {
        Type: type,
        Source: clientId,
        Destination: destinationId,
        Text: text,
        Session: session
    };
    conn.value.send(JSON.stringify(msg));
    return session; // Return the session ID for tracking
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
  
  return { conn, wsUrl, connect, sendMessage }
});