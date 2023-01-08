<template>
  <div class="card m-3">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <h4 class="card-header">User</h4>

    <h4 class="card-header">Create new user</h4>
    <div class="card-body">
      <div class="form-group">
        <form @submit.prevent="createUser">
          <label for="email">Email</label>
          <input type="email" v-model="email" class="form-control" />
          <br />
          <label for="password">Password</label>
          <input type="password" v-model="password" class="form-control" />
          <br />
          <label for="roles">Roles</label>
          <select v-model="selectedRoles" multiple class="form-control">
            <option v-for="role in roles" :key="role" :value="role">{{ role }}</option>
          </select>
          <br />
          <button class="btn btn-primary pointer" type="submit">Create User</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import http from '@/common-http';

export default defineComponent({
  name: 'UserAdmin',
  data() {
    return {
      email: '',
      password: '',
      selectedRoles: [],
      roles: ['Admin', 'Maintainer', 'Tester'],
      error: '',
      loading: false,
    };
  },
  methods: {
    async createUser() {
      this.loading = true;
      this.error = '';

      const newUser = {
        email: this.email,
        password: this.password,
        roles: this.selectedRoles,
      };

      http.post('/api/v1/users', newUser).catch((err) => {
        this.error = err + ' | ' + err.response?.data?.error;
      });
      this.loading = false;
      this.email = '';
      this.password = '';
      this.selectedRoles = [];
    },
  },
});
</script>
