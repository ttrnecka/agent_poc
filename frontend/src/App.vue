<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { onMounted, computed, watch } from 'vue'
import Header from './components/Header.vue'
import { useDataStore } from '@/stores/dataStore'
import ws from '@/services/websocket'
import router from '@/router'

const dataStore = useDataStore()

watch(() => dataStore.isLoggedIn, (val) => {
  if (val) {
    dataStore.getData()
    ws.connectWebSocket()
  } else {
    router.push('/login')
  }
})

onMounted(() => {
  dataStore.fetchUser();
})

</script>

<template>
  <div>
    <header>
      <img alt="Vue logo" class="logo" src="@/assets/xa.png" width="200" height="125" />

      <div class="wrapper">
        <Header msg="Welcome to xAnalytics" />

        <nav>
          <RouterLink to="/">Configuration</RouterLink>
          <RouterLink to="/collectors">Collector</RouterLink>
          <RouterLink to="/inventory">Inventory</RouterLink>
          <button
              v-if="dataStore.isLoggedIn"
              @click="dataStore.logoutUser"
              class="btn btn-outline-secondary"
            >
              Logout
            </button>
          </nav>
      </div>
    </header>

    <RouterView />
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: stretch;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
