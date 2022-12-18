<template>
  <nav class="navbar navbar-expand-lg">
    <div class="container-fluid">
      <a class="navbar-brand" href="/"> <img src="./assets/logo.png" alt="Test and Win - e2e test coverage" width="110" height="50" /></a>
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarSupportedContent"
        aria-controls="navbarSupportedContent"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li v-if="signedin && isEditor" class="nav-item">
            <a class="nav-link" href="/product">Product</a>
          </li>
          <li v-if="signedin && isConsumer" class="nav-item">
            <a class="nav-link" href="/coverage">Coverage</a>
          </li>
          <li v-if="signedin && isConsumer" class="nav-item">
            <a class="nav-link" href="/tests">Tests</a>
          </li>
          <li v-if="!signedin" class="nav-item">
            <a class="nav-link" href="/signin">Sign In</a>
          </li>
          <li v-if="signedin" class="nav-item">
            <a class="nav-link" href="/signout">Sign Out</a>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <router-view />
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'App',
  data() {
    return {
      signedin: sessionStorage.getItem('token') != undefined,
    };
  },
  computed: {
    isEditor() {
      const s = sessionStorage.getItem('roles');
      if (s) {
        return s.indexOf('e') > -1;
      }
      return false;
    },
    isConsumer() {
      const s = sessionStorage.getItem('token');
      if (s) {
        return s.indexOf('c') > -1;
      }
      return false;
    }
  }
})
</script>

<style>
@import "./assets/styles.css";
</style>
