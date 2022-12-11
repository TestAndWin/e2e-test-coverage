<template>
  <div class="">
    <SignIn />
  </div>
</template>

<script lang="ts">
import SignIn from "@/components/SignIn.vue";

export default {
  components: {
    SignIn,
  },
  mounted() {
    document.title = "e2e coverage - Sign In";
  },
};
</script>
