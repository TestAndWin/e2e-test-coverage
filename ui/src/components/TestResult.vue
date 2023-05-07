<template>
  <span v-if="props.test['total'] < 1" class="result failures">
    {{ props.test['total'] }}
  </span>
  <span v-if="props.test['total'] > 0" class="result total">
    {{ props.test['total'] }}
    <i v-if="props.test['total'] > props.test['first-total']" class="bi bi-caret-up"></i>
    <i v-if="props.test['total'] < props.test['first-total']" class="bi bi-caret-down"></i>
  </span>
  &nbsp;
  <span v-if="props.test['passes'] > 0" class="result passes">
    {{ props.test['passes'] }}
  </span>
  <span v-if="props.test['passes'] < 1" class="result">-</span>
  &nbsp;
  <span v-if="props.test['failures'] > 0" class="result failures">
    {{ props.test['failures'] }}
  </span>
  <span v-if="props.test['failures'] < 1" class="result">-</span>
  &nbsp;
  <span v-if="props.test['pending'] > 0" class="result pending">
    {{ props.test['pending'] }}
  </span>
  <span v-if="props.test['pending'] < 1" class="result">-</span>
  &nbsp;
  <span v-if="props.test['skipped'] > 0" class="result skipped">
    {{ props.test['skipped'] }}
  </span>
  <span v-if="props.test['skipped'] < 1" class="result">-</span>
</template>

<script setup lang="ts">
const props = defineProps({
  test: { type: Object, required: true },
});
</script>
