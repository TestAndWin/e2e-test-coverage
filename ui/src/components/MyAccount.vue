<template>
  <div class="product container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="">Change Password</h4>
    <div class="user shadow p-2 mb-4 rounded">
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

    <div v-if="isAdmin()">
      <h4 class="">Generate API Key</h4>
      <div class="user shadow p-2 mb-4 rounded">
        <div class="form-group">
          <button class="btn btn-primary pointer" @click="generateApiKey()">Generate</button>
        </div>
        <br />
        <div class="form-group">
          <label>API Key:</label>&nbsp;
          <span>{{ apiKey }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import http from '@/common-http';
import { isAdmin } from '@/menu';
const router = useRouter();

const loading = ref(false);
const error = ref('');

const apiKey = ref('');
const generateApiKey = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await http.post(`/api/v1/users/generate-api-key`, {});

    // Extract data from StandardResponse format
    if (response.data && response.data.data && response.data.data.key) {
      apiKey.value = response.data.data.key;
    } else {
      error.value = 'Could not extract API key from response';
    }
  } catch (err) {
    error.value = `Error generating API key: ${err}`;
  } finally {
    loading.value = false;
  }
};

const newPassword = ref('');
const password = ref('');
const changePassword = async () => {
  loading.value = true;
  error.value = '';

  http
    .put(`/api/v1/users/change-pwd`, {
      'new-password': newPassword.value,
      password: password.value
    })
    .then(() => {
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
