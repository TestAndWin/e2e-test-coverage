<template>
  <div class="product container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="">User</h4>
    <div id="user" class="user shadow p-2 mb-4 rounded">
      <div v-for="u in user" :key="u['id']">
        <div class="row">
          <div class="col">
            <span class="justify-content-between pointer">{{ u['email'] }}</span>
            {{ u.roles }}
            <a @click="showEditUserModal(u.id ?? 0, u.email ?? '', u.roles ?? [])"
              ><i class="bi bi-pencil pointer"></i></a
            >&nbsp;
            <a @click="deleteUser(u.id ?? 0)"><i class="bi bi-trash pointer"></i></a>
          </div>
        </div>
        <hr />
      </div>
    </div>
    <button class="btn btn-primary pointer" type="submit" @click="showAddUserModal()">Add User</button>

    <!-- Modal to add an user -->
    <div class="modal fade" id="addUser" tabindex="-1" aria-labelledby="addUserLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="addUserLabel">Add User</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <form>
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
              <button class="btn btn-primary pointer" type="submit" data-bs-dismiss="modal" @click="createUser">
                Add
              </button>
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Modal to edit an user -->
    <div class="modal fade" id="editUser" tabindex="-1" aria-labelledby="editUserLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="editUserLabel">Edit User</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <form>
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
              <button type="submit" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="saveUser">
                Save
              </button>
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import http from '@/common-http';
import { Modal } from 'bootstrap';
import type { User } from '@/types';

const roles = ref(['Admin', 'Maintainer', 'Tester']);
const error = ref('');
const loading = ref(false);
const email = ref('');
const password = ref('');
const selectedRoles = ref(['']);
const user = ref<User[]>([]);

const showAddUserModal = () => {
  email.value = '';
  password.value = '';
  selectedRoles.value = [];
  new Modal('#addUser').show();
};

const getUser = async () => {
  loading.value = true;
  try {
    const response = await http.get(`/api/v1/users`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      user.value = response.data.data;
    } else {
      user.value = [];
    }
  } catch (err) {
    error.value = `Error loading users: ${err}`;
    user.value = []; // Initialize as empty array on error
  } finally {
    loading.value = false;
  }
};

const userId = ref(0);
const showEditUserModal = (uId: number, e: string, r: string | string[]) => {
  userId.value = uId;
  email.value = e;

  // Handle roles based on type
  if (Array.isArray(r)) {
    // If it's already an array, use it directly
    selectedRoles.value = r;
  } else if (typeof r === 'string') {
    // If it's a comma-separated string, split it
    selectedRoles.value = r.split(',').map((role) => role.trim());
  } else {
    // Default to empty if undefined or other type
    selectedRoles.value = [];
  }

  console.log('Selected roles set to:', selectedRoles.value);

  // Show modal
  new Modal('#editUser').show();
};

const deleteUser = async (userId: number) => {
  try {
    loading.value = true;

    await http.delete(`/api/v1/users/${userId}`);

    // Refresh user list
    await getUser();
  } catch (err) {
    error.value = `Error deleting user: ${err}`;
  } finally {
    loading.value = false;
  }
};

const createUser = async () => {
  loading.value = true;
  error.value = '';

  try {
    // Validate inputs
    if (!email.value || !password.value || selectedRoles.value.length === 0) {
      error.value = 'Email, password, and at least one role are required';
      return;
    }

    const newUser = {
      email: email.value,
      password: password.value,
      roles: selectedRoles.value
    };

    await http.post('/api/v1/users', newUser);

    // Clear form
    email.value = '';
    password.value = '';
    selectedRoles.value = [];

    // Refresh user list
    await getUser();
  } catch (err) {
    error.value = `Error creating user: ${err}`;
  } finally {
    loading.value = false;
  }
};

const saveUser = async () => {
  loading.value = true;
  error.value = '';

  try {
    // Validate inputs
    if (!email.value || selectedRoles.value.length === 0) {
      error.value = 'Email and at least one role are required';
      return;
    }

    // Prepare user data (password can be empty for updates)
    const userData = {
      email: email.value,
      password: password.value,
      roles: selectedRoles.value
    };

    await http.put(`/api/v1/users/${userId.value}`, userData);

    // Clear form
    userId.value = 0;
    email.value = '';
    password.value = '';
    selectedRoles.value = [];

    // Refresh user list
    await getUser();
  } catch (err) {
    error.value = `Error updating user: ${err}`;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  getUser();
});
</script>
