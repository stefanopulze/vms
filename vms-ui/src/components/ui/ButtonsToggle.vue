<script setup lang="ts">
defineProps<{ values: { label: string, value: string | number }[] }>()
const model = defineModel<string | number>()
const emits = defineEmits(['change'])

function setValue(v: string | number) {
  const oldValue = model.value
  if (v != oldValue) {
    model.value = v
    emits('change', v, oldValue)
  }
}

</script>

<template>
  <div class="buttons toggle">
    <button v-for="v of values"
            :key="v.value"
            :class="{active: model === v.value}"
            @click.prevent="setValue(v.value)">{{ v.label }}
    </button>
  </div>
</template>


