import { getClientId } from './auth.js'
import ws from '@/services/websocket'
import { useApiStore } from '@/stores/apiStore';
import { v4 as uuidv4 } from 'uuid';

export const MESSAGE_TYPE = {
  'COLLECTOR_ONLINE': 1,
  'COLLECTOR_OFFLINE': 2,
  'POLICY_REFRESH': 10,
  'PROBE_START': 20,
  'PROBE_RUNNING': 21,
  'PROBE_FINISHED_OK': 22,
  'PROBE_FINISHED_ERR': 23,
  'PROBE_DATA': 24
};

export function sendMessage(type, destinationId, text, session = uuidv4()) {
    const clientId = getClientId(); 
    const msg = {
        Type: type,
        Source: clientId,
        Destination: destinationId,
        Text: text,
        Session: session
    };
    ws.getConn().send(JSON.stringify(msg));
    return session; // Return the session ID for tracking
  }

export function handleMessage(msg) {
  const apiStore = useApiStore();
  
  console.log('📥 Message:', msg);
  if (msg.Source != getClientId()) {
    if (msg.Type == MESSAGE_TYPE.COLLECTOR_ONLINE) {
      apiStore.updateCollectorStatus(msg.Source,"ONLINE")
    }
    if (msg.Type == MESSAGE_TYPE.COLLECTOR_OFFLINE) {
      apiStore.updateCollectorStatus(msg.Source,"OFFLINE")
    }
  }
}
