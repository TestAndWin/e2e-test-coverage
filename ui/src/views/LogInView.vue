<template>
  <div class="">
    <LogIn />
  </div>
</template>

<script setup lang="ts">
import LogIn from "@/components/LogIn.vue";
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - Log in';
});
</script>
