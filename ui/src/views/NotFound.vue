<template>
  <div class="container">
    <div class="card">
      <div class="card-body">
        <h1>Not Found</h1>
        <p>Oops, we couldn't find that page.</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  mounted() {
    document.title = "e2e coverage - 404";
  },
};
</script>
