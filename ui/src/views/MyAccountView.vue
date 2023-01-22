<template>
  <div class="">
    <MyAccount />
  </div>
</template>

<script lang="ts">
import MyAccount from "@/components/MyAccount.vue";

export default {
  components: {
    MyAccount,
  },
  mounted() {
    document.title = "e2e coverage - My Account";
  },
};
</script>
