<template>
  <input
    v-bind="$attrs"
    :value="modelValue"
    @input="onInput"
    :class="classes"
  />
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: "",
  },
})

const emit = defineEmits(["update:modelValue"])

const onInput = (event) => {
  emit("update:modelValue", event.target.value)
}

const classes =
  "mt-2 w-full rounded-2xl border border-border bg-white/90 px-4 py-2 text-sm text-foreground shadow-sm focus:border-primary focus:outline-none focus:ring-2 focus:ring-ring/30"
</script>