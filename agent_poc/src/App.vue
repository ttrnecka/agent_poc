<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { onMounted } from 'vue'
import Header from './components/Header.vue'
import { useApiStore } from '@/stores/apiStore'
import { useWsConnectionStore } from '@/stores/wsStore'

const apiStore = useApiStore()
const ws = useWsConnectionStore()

onMounted(() => {
  apiStore.loadCollectors()
  apiStore.loadPolicies()
  apiStore.loadProbes()
  ws.connect()
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
