<template>
  <div class="">
    <UserAdmin />
  </div>
</template>

<script lang="ts">
import UserAdmin from "@/components/UserAdmin.vue";

export default {
  components: {
    UserAdmin,
  },
  mounted() {
    document.title = "e2e coverage - Admin";
  },
};
</script>
