<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">
          All tests of the product with their status of the last 28 days. Also tests not assigned to an area or feature are displayed. Only the latest test per suite name will be
          displayed.
        </p>
      </div>
    </div>
  </div>
  <div class="">
    <TestCoverage :productId="1" />
  </div>
</template>

<script lang="ts">
import TestCoverage from "@/components/TestCoverage.vue";

export default {
  components: {
    TestCoverage,
  },
  mounted() {
    document.title = "e2e coverage - Coverage";
  },
};
</script>
