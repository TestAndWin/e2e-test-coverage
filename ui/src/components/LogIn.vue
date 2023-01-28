<template>
  <div class="card m-3">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" class="spinner-border info" role="status">
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

<script setup lang="ts">
import { ref } from 'vue';
//import { useRouter } from 'vue-router';
import http from '@/common-http';
//const router = useRouter();

const loading = ref(false);
const error = ref('');

const email = ref('admin');
const password = ref('e2ecoverage');
const login = async () => {
  loading.value = true;
  error.value = '';

  http
    .post('/api/v1/auth/login', { email: email.value, password: password.value })
    .then((response) => {
      loading.value = false;
      sessionStorage.setItem('token', response.data.token);
      sessionStorage.setItem('refreshToken', response.data.refreshToken);
      sessionStorage.setItem('roles', response.data.roles);
      location.assign('/');
      //router.push('/');
    })
    .catch(() => {
      error.value = 'Log in failed';
    })
    .finally(() => {
      email.value = '';
      password.value = '';
      loading.value = false;
    });
};
</script>

<style scoped>
@import '../assets/styles.css';
</style>
