<template>
  <div class="card m-3">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="card-header">Log In</h4>
    <div class="card-body">
      <div class="form-group">
        <label>E-Mail</label>
        <input v-model="email" type="text" class="form-control" />
      </div>
      <div class="form-group">
        <label>Password</label>
        <input v-model="password" type="password" class="form-control" />
      </div>
      <br />
      <div class="form-group">
        <button class="btn btn-primary pointer" @click="login()">Log In</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import http from '@/common-http';

export default defineComponent({
  name: 'LogIn',
  data() {
    return {
      email: 'admin',
      password: 'e2ecoverage',
      loading: false,
      error: '',
    };
  },
  methods: {
    async login() {
      this.loading = true;
      this.error = '';

      http
        .post('/api/v1/auth/login', { email: this.email, password: this.password })
        .then((response) => {
          this.loading = false;
          sessionStorage.setItem('token', response.data.token);
          sessionStorage.setItem('refreshToken', response.data.refreshToken);
          sessionStorage.setItem('roles', response.data.roles);
          location.assign('/');
        })
        .catch(() => {
          this.error = 'Log in failed';
          this.email = '';
          this.password = '';
          this.loading = false;
        });
    },
  },
  components: {},
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
