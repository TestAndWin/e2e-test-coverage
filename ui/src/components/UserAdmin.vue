<template>
  <div class="product container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="">User</h4>
    <div class="user shadow p-2 mb-4 rounded" :key="componentKey">
      <div v-for="u in user" :key="u['id']">
        <div class="row">
          <div class="col">
            <span class="justify-content-between pointer">{{ u['email'] }}</span>
            {{ u['roles'] }}
            <a @click="showEditUserModal(u['id'], u['email'], u['roles'])"><i class="bi bi-pencil pointer"></i></a>&nbsp;
            <a @click="deleteUser(u['id'])"><i class="bi bi-trash pointer"></i></a>
          </div>
        </div>
        <hr />
      </div>
    </div>
    <button class="btn btn-primary pointer" type="submit" @click="showAddUserModal()">Add User</button>

    <!-- Modal to add an user -->
    <div class="modal fade" ref="addUser" id="addUser" tabindex="-1" aria-labelledby="addUserLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="addUserLabel">Add User</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <form @submit.prevent="createUser">
            <div class="modal-body">
              <div class="form-group">
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
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-primary pointer" type="submit" data-bs-dismiss="modal">Add</button>
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Modal to edit an user -->
    <div class="modal fade" ref="editUser" id="editUser" tabindex="-1" aria-labelledby="editUserLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="editUserLabel">Edit User</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <form @submit.prevent="saveUser">
            <div class="modal-body">
              <div class="form-group">
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
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-primary pointer" type="submit" data-bs-dismiss="modal">Save</button>
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>    
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import http from '@/common-http';
import { Modal } from 'bootstrap';

export default defineComponent({
  name: 'UserAdmin',
  data() {
    return {
      email: '',
      password: '',
      selectedRoles: [''],
      roles: ['Admin', 'Maintainer', 'Tester'],
      error: '',
      user: [],
      userId: 0,
      loading: false,
      componentKey: 0,
    };
  },
  methods: {
    forceRerender() {
      this.componentKey += 1;
    },
    async showAddUserModal() {
      new Modal('#addUser').show();
    },
    async showEditUserModal(userId: number, email: string, roles: string) {
      this.userId = userId;
      this.email = email;
      this.selectedRoles = Array.from(roles);
      new Modal('#editUser').show();
    },
    async getUser() {
      this.loading = true;
      await http
        .get(`/api/v1/users`)
        .then((response) => {
          this.user = response.data;
        })
        .catch((err) => {
          this.error = err + ' | ' + err.response.data.error;
        });
      this.loading = false;
      this.forceRerender();
    },
    async deleteUser(userId: number) {
      await http.delete(`/api/v1/users/${userId}`).catch((err) => {
        this.error = err + ' | ' + err.response.data.error;
      });
      this.getUser();
    },
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
      this.getUser();
    },
    async saveUser() {
      this.loading = true;
      this.error = '';

      const newUser = {
        email: this.email,
        password: this.password,
        roles: this.selectedRoles,
      };

      http.put(`/api/v1/users/${this.userId}`, newUser).catch((err) => {
        this.error = err + ' | ' + err.response?.data?.error;
      });
      this.userId = 0;
      this.loading = false;
      this.email = '';
      this.password = '';
      this.selectedRoles = [];
      this.getUser();
    },
  },
  mounted() {
    this.getUser();
  },
});
</script>
