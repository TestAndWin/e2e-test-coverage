<template>
  <div class="">
    <LogIn />
  </div>
</template>

<script lang="ts">
import LogIn from "@/components/LogIn.vue";

export default {
  components: {
    LogIn,
  },
  mounted() {
    document.title = "e2e coverage - Log In";
  },
};
</script>
