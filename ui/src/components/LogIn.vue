<template>
  <div class="product container">
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
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import http from '@/common-http';

// Helper function to safely store user roles in multiple places for redundancy
function storeUserRoles(roles: string) {
  try {
    // Store in sessionStorage (primary storage)
    sessionStorage.setItem('roles', roles);

    // Store in localStorage as backup (will persist between sessions)
    localStorage.setItem('roles_backup', roles);

    console.log('Roles stored successfully in both session and local storage');
  } catch (error) {
    console.error('Error storing roles:', error);
  }
}

const loading = ref(false);
const error = ref('');

const email = ref('');
const password = ref('');
const login = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await http.post('/api/v1/auth/login', {
      email: email.value,
      password: password.value
    });

    console.log('Login response:', response.data);

    // Trace the entire response structure for debugging
    const traceObject = (obj: any, path = '') => {
      if (!obj || typeof obj !== 'object') return;

      Object.entries(obj).forEach(([key, value]) => {
        const currentPath = path ? `${path}.${key}` : key;
        console.log(`${currentPath}:`, value);

        if (typeof value === 'object' && value !== null) {
          traceObject(value, currentPath);
        }
      });
    };

    console.log('FULL RESPONSE TRACE:');
    traceObject(response.data);

    // Try all possible locations where roles might be
    let rolesValue = null;

    // Direct in response data
    if (response.data && response.data.roles) {
      rolesValue = response.data.roles;
    }
    // In data property of StandardResponse
    else if (response.data && response.data.data && response.data.data.roles) {
      rolesValue = response.data.data.roles;
    }
    // Try checking success property first
    else if (response.data && response.data.success === true) {
      const userData = response.data.data;
      if (userData && userData.roles) {
        rolesValue = userData.roles;
      }
    }

    if (rolesValue) {
      console.log('Found roles in response:', rolesValue);

      // Store roles in both sessionStorage and localStorage
      storeUserRoles(rolesValue);

      // Also store userId if available
      if (response.data && response.data.data && response.data.data.userId) {
        console.log('Storing user ID:', response.data.data.userId);
        localStorage.setItem('userId', response.data.data.userId);
      }

      // Log to verify it was set correctly
      console.log('Roles after setting:', sessionStorage.getItem('roles'));

      // Redirect to home page to refresh the menu
      location.assign('/');
      return;
    }

    // If we get here, we couldn't find the roles
    console.error('Login response missing roles:', response.data);
    error.value = 'Login successful, but the system returned incomplete data. Please try again.';
  } catch (err) {
    console.error('Login error:', err);
    error.value = 'Login failed. Please check your credentials and try again.';
  } finally {
    // Clear form fields and loading state
    email.value = '';
    password.value = '';
    loading.value = false;
  }
};
</script>

<style scoped>
@import '../assets/styles.css';
</style>
