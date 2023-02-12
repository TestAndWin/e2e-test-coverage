<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Maintaining product areas and their features.</p>
      </div>
    </div>
  </div>
  <div class="">
    <ProductComp :productId=1 />
  </div>
</template>

<script setup lang="ts">
import ProductComp from "@/components/ProductComp.vue";
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - Product';
});
</script>
