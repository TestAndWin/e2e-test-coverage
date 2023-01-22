<template>
  <div class="card m-3">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
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

<script lang="ts">
import { defineComponent } from 'vue';
import http from '@/common-http';

export default defineComponent({
  name: 'LogIn',
  data() {
    return {
      newPassword: '',
      password: '',
      loading: false,
      error: '',
    };
  },
  methods: {
    async changePassword() {
      this.loading = true;
      this.error = '';

      http
        .put(`/api/v1/users/change-pwd`, { "new-password": this.newPassword, password: this.password })
        .then(() => {
          this.loading = false;
          location.assign('/');
        })
        .catch((err) => {
          this.error = err + ' | ' + err.response?.data?.error;
        });
      this.newPassword = '';
      this.password = '';
      this.loading = false;
    },
  },
  components: {},
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
