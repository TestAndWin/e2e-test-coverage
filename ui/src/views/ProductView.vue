<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Maintain products areas and their features.</p>
      </div>
    </div>
  </div>
  <div class="">
    <ProductComp :productId="1" />
  </div>
</template>

<script lang="ts">
import ProductComp from "@/components/ProductComp.vue";

export default {
  components: {
    ProductComp,
  },
  mounted() {
    document.title = "e2e coverage - Product";
  },
};
</script>
