<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import {
  CheckCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  TrashIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import { appDialogQueue, resolveAppDialog } from '@/composables/useAppDialog'

const dialog = computed(() => appDialogQueue.value[0] || null)
const cancelButton = ref<HTMLButtonElement | null>(null)
const confirmButton = ref<HTMLButtonElement | null>(null)

const iconByTone = {
  success: CheckCircleIcon,
  error: ExclamationTriangleIcon,
  warning: ExclamationTriangleIcon,
  info: InformationCircleIcon,
  danger: TrashIcon
}

const toneClasses = {
  success: 'bg-emerald-100 text-emerald-600 ring-emerald-100 dark:bg-emerald-500/15 dark:text-emerald-300 dark:ring-emerald-500/20',
  error: 'bg-rose-100 text-rose-600 ring-rose-100 dark:bg-rose-500/15 dark:text-rose-300 dark:ring-rose-500/20',
  warning: 'bg-amber-100 text-amber-600 ring-amber-100 dark:bg-amber-500/15 dark:text-amber-300 dark:ring-amber-500/20',
  info: 'bg-indigo-100 text-indigo-600 ring-indigo-100 dark:bg-indigo-500/15 dark:text-indigo-300 dark:ring-indigo-500/20',
  danger: 'bg-rose-100 text-rose-600 ring-rose-100 dark:bg-rose-500/15 dark:text-rose-300 dark:ring-rose-500/20'
}

const confirmButtonClasses = computed(() =>
  dialog.value?.tone === 'danger'
    ? 'bg-rose-600 text-white shadow-rose-600/20 hover:bg-rose-700 focus-visible:ring-rose-500'
    : 'bg-indigo-600 text-white shadow-indigo-600/20 hover:bg-indigo-700 focus-visible:ring-indigo-500'
)

const dismiss = (confirmed: boolean): void => {
  resolveAppDialog(confirmed)
}

const onKeydown = (event: KeyboardEvent): void => {
  if (event.key === 'Escape' && dialog.value) {
    event.preventDefault()
    dismiss(false)
  }
}

watch(dialog, async currentDialog => {
  if (!currentDialog) return
  await nextTick()
  if (currentDialog.kind === 'confirm') cancelButton.value?.focus()
  else confirmButton.value?.focus()
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="dialog"
        class="fixed inset-0 z-[300] flex items-center justify-center overflow-y-auto bg-slate-950/60 p-4 backdrop-blur-sm"
        @keydown="onKeydown"
        @click.self="dialog.kind === 'alert' && dismiss(true)"
      >
        <Transition
          appear
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="translate-y-3 scale-95 opacity-0"
          enter-to-class="translate-y-0 scale-100 opacity-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="translate-y-0 scale-100 opacity-100"
          leave-to-class="translate-y-2 scale-95 opacity-0"
        >
          <section
            v-if="dialog"
            :key="dialog.id"
            role="alertdialog"
            aria-modal="true"
            :aria-labelledby="`app-dialog-title-${dialog.id}`"
            :aria-describedby="`app-dialog-message-${dialog.id}`"
            class="w-full max-w-md overflow-hidden rounded-3xl border border-white/60 bg-white shadow-2xl shadow-slate-950/25 dark:border-slate-700 dark:bg-slate-900"
          >
            <div class="relative px-6 pb-6 pt-7 text-center sm:px-8">
              <button
                v-if="dialog.kind === 'alert'"
                type="button"
                class="absolute right-4 top-4 rounded-xl p-2 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 dark:hover:bg-slate-800 dark:hover:text-slate-200"
                aria-label="Tutup"
                @click="dismiss(true)"
              >
                <XMarkIcon class="h-5 w-5" />
              </button>

              <div
                class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl ring-1"
                :class="toneClasses[dialog.tone]"
              >
                <component :is="iconByTone[dialog.tone]" class="h-7 w-7" />
              </div>

              <h2
                :id="`app-dialog-title-${dialog.id}`"
                class="text-xl font-extrabold tracking-tight text-slate-900 dark:text-white"
              >
                {{ dialog.title }}
              </h2>
              <p
                :id="`app-dialog-message-${dialog.id}`"
                class="mx-auto mt-2 max-w-sm whitespace-pre-line text-sm leading-6 text-slate-600 dark:text-slate-300"
              >
                {{ dialog.message }}
              </p>
            </div>

            <div class="flex flex-col-reverse gap-3 border-t border-slate-100 bg-slate-50/80 px-6 py-4 dark:border-slate-800 dark:bg-slate-950/50 sm:flex-row sm:justify-end sm:px-8">
              <button
                v-if="dialog.kind === 'confirm'"
                ref="cancelButton"
                type="button"
                class="inline-flex min-h-11 items-center justify-center rounded-xl border border-slate-200 bg-white px-5 text-sm font-bold text-slate-700 transition hover:bg-slate-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-slate-700"
                @click="dismiss(false)"
              >
                {{ dialog.cancelLabel }}
              </button>
              <button
                ref="confirmButton"
                type="button"
                class="inline-flex min-h-11 items-center justify-center rounded-xl px-5 text-sm font-bold shadow-lg transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 dark:focus-visible:ring-offset-slate-900"
                :class="confirmButtonClasses"
                @click="dismiss(true)"
              >
                {{ dialog.confirmLabel }}
              </button>
            </div>
          </section>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
