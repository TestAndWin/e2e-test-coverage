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
              <button class="btn btn-primary pointer" type="submit" data-bs-dismiss="modal" @click="createUser">Add</button>
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
              <button type="submit" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="saveUser">Save</button>
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

const roles = ref(['Admin', 'Maintainer', 'Tester']);
const error = ref('');
const loading = ref(false);
const email = ref('');
const password = ref('');
const selectedRoles = ref(['']);
const user = ref([]);

const showAddUserModal = () => {
  email.value = '';
  password.value = '';
  selectedRoles.value = [];
  new Modal('#addUser').show();
};

const getUser = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/users`)
    .then((response) => {
      user.value = response.data;
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response.data.error;
    });
  loading.value = false;
};

const userId = ref(0);
const showEditUserModal = (uId: number, e: string, r: string) => {
  userId.value = uId;
  email.value = e;
  selectedRoles.value = Array.from(r);
  new Modal('#editUser').show();
};

const deleteUser = async (userId: number) => {
  await http.delete(`/api/v1/users/${userId}`).catch((err) => {
    error.value = err + ' | ' + err.response.data.error;
  });
  getUser();
};

const createUser = async () => {
  loading.value = true;
  error.value = '';

  const newUser = {
    email: email.value,
    password: password.value,
    roles: selectedRoles.value,
  };

  await http
    .post('/api/v1/users', newUser)
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    })
    .finally(() => {
      loading.value = false;
      email.value = '';
      password.value = '';
      selectedRoles.value = [];
    });
  await getUser();
};

const saveUser = async () => {
  loading.value = true;
  error.value = '';

  const newUser = {
    email: email.value,
    password: password.value,
    roles: selectedRoles.value,
  };
  console.log('aaveUser');
  console.log(newUser);

  await http
    .put(`/api/v1/users/${userId.value}`, newUser)
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    })
    .finally(() => {
      userId.value = 0;
      loading.value = false;
      email.value = '';
      password.value = '';
      selectedRoles.value = [];
    });
  await getUser();
};

onMounted(() => {
  getUser();
});
</script>
