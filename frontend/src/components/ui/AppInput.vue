<script setup lang="ts">
defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  placeholder: {
    type: String,
    default: ''
  },
  error: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="flex flex-col gap-1.5 w-full">
    <!-- Teks Label di Atas Input -->
    <label v-if="label" class="text-sm font-semibold text-text-secondary dark:text-slate-300">
      {{ label }} <span v-if="required" class="text-accent-danger">*</span>
    </label>
    
    <!-- Input Box (100% Mengkopi Desain .form-control Bapak) -->
    <input
      :type="type"
      :placeholder="placeholder"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :required="required"
      :class="[
        'w-full px-3.5 py-2.5 text-[0.9rem] rounded-[10px] transition-all duration-150 ease-in-out focus:outline-none focus:ring-[3px]',
        'bg-bg-secondary text-tx-primary border',
        'dark:bg-slate-900 dark:border-slate-700 dark:text-slate-50',
        error 
          ? 'border-accent-danger focus:border-accent-danger focus:ring-accent-danger/20' 
          : 'border-border-color focus:border-accent-primary focus:ring-accent-primary/12'
      ]"
    />
    
    <!-- Teks Error Merah di Bawah (Jika Validasi Gagal) -->
    <span v-if="error" class="text-xs text-accent-danger font-medium">{{ error }}</span>
  </div>
</template>
