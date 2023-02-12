<template>
    <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Change your password and user with the Admin role have the ability to generate an API token.</p>
      </div>
    </div>
  </div>
  <div class="">
    <MyAccount />
  </div>
</template>

<script setup lang="ts">
import MyAccount from "@/components/MyAccount.vue";
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - My Account';
});
</script>
