<template>
  <div class="card m-3">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="card-header">Change Password</h4>
    <div class="card-body">
      <div class="form-group">
        <label>Password</label>
        <input v-model="password" type="password" class="form-control" />
      </div>
      <div class="form-group">
        <label>New Password</label>
        <input v-model="newPassword" type="password" class="form-control" />
      </div>
      <br />
      <div class="form-group">
        <button class="btn btn-primary pointer" @click="changePassword()">Save</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import http from '@/common-http';
const router = useRouter();

const loading = ref(false);
const error = ref('');

const newPassword = ref('');
const password = ref('');
const changePassword = async () => {
  loading.value = true;
  error.value = '';

  http
    .put(`/api/v1/users/change-pwd`, { 'new-password': newPassword.value, password: password.value })
    .then(() => {
      loading.value = false;
      router.push('/');
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    })
    .finally(() => {
      newPassword.value = '';
      password.value = '';
      loading.value = false;
    });
};
</script>

<style scoped>
@import '../assets/styles.css';
</style>
