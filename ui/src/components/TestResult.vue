<template>
  <span v-if="total < 1" class="result failures">
    {{ total }}
  </span>
  <span v-if="total > 0" class="result total">
    {{ total }}
    <i v-if="total > firstTotal" class="bi bi-caret-up"></i>
    <i v-if="total < firstTotal" class="bi bi-caret-down"></i>
  </span>
  &nbsp;
  <span v-if="passes > 0" class="result passes">
    {{ passes }}
  </span>
  <span v-if="passes < 1" class="result">-</span>
  &nbsp;
  <span v-if="failures > 0" class="result failures">
    {{ failures }}
  </span>
  <span v-if="failures < 1" class="result">-</span>
  &nbsp;
  <span v-if="pending > 0" class="result pending">
    {{ pending }}
  </span>
  <span v-if="pending < 1" class="result">-</span>
  &nbsp;
  <span v-if="skipped > 0" class="result skipped">
    {{ skipped }}
  </span>
  <span v-if="skipped < 1" class="result">-</span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps({
  test: { type: Object, required: true }
});

// Compute values with fallbacks to prevent NaN
const total = computed(() => {
  const value = props.test['total'] || 0;
  return isNaN(value) ? 0 : value;
});

const firstTotal = computed(() => {
  const value = props.test['first-total'] || 0;
  return isNaN(value) ? 0 : value;
});

const passes = computed(() => {
  const value = props.test['passes'] || 0;
  return isNaN(value) ? 0 : value;
});

const failures = computed(() => {
  const value = props.test['failures'] || 0;
  return isNaN(value) ? 0 : value;
});

const pending = computed(() => {
  const value = props.test['pending'] || 0;
  return isNaN(value) ? 0 : value;
});

const skipped = computed(() => {
  const value = props.test['skipped'] || 0;
  return isNaN(value) ? 0 : value;
});
</script>
