import { v4 as uuidv4 } from 'uuid';
import { ref } from 'vue'

export function getClientId() {
  let clientId = localStorage.getItem("client_id");
  if (!clientId) {
    clientId = "web_client_" + uuidv4(); // Secure, random
    localStorage.setItem("client_id", clientId);
  }
  return clientId;
}

// // reactive login state
// export const loggedIn = ref(!!localStorage.getItem('user'))

// export function isLoggedIn() {
//   return loggedIn.value
// }

// export function setLoggedIn(value) {
//   localStorage.setItem('user', value)
//   loggedIn.value = value
//   if (!value) {
//     localStorage.removeItem('user')
//   }
// }
