<template>
  <div class=""></div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
//import { useRouter } from 'vue-router';

//const router = useRouter()
onMounted(() => {
  sessionStorage.clear();
  //router.push('/');
  location.assign('/');
});
</script>
