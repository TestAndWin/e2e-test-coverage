<template>
  <div class=""></div>
</template>

<script lang="ts">
export default {
  mounted() {
    sessionStorage.clear();
    location.assign('/');
  },
};
</script>
