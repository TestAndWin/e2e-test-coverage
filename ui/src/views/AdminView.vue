<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Maintain user.</p>
      </div>
    </div>
  </div>
  <div class="">
    <UserAdmin />
  </div>
</template>

<script lang="ts">
import UserAdmin from '@/components/UserAdmin.vue';

export default {
  components: {
    UserAdmin,
  },
  mounted() {
    document.title = 'e2e coverage - Admin';
  },
};
</script>
