<template>
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
