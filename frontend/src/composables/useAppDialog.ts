import { ref } from 'vue'

export type DialogTone = 'success' | 'error' | 'warning' | 'info' | 'danger'

export interface DialogRequest {
  id: number
  kind: 'alert' | 'confirm'
  title: string
  message: string
  tone: DialogTone
  confirmLabel: string
  cancelLabel: string
  resolve: (confirmed: boolean) => void
}

const dialogQueue = ref<DialogRequest[]>([])
let nextDialogId = 0

const enqueueDialog = (dialog: Omit<DialogRequest, 'id' | 'resolve'>): Promise<boolean> =>
  new Promise(resolve => {
    dialogQueue.value.push({
      ...dialog,
      id: ++nextDialogId,
      resolve
    })
  })

export const showAppAlert = async (
  message: string,
  tone: Exclude<DialogTone, 'danger'> = 'info',
  title?: string
): Promise<void> => {
  const titles = {
    success: 'Berhasil',
    error: 'Terjadi Kesalahan',
    warning: 'Perlu Perhatian',
    info: 'Informasi'
  }

  await enqueueDialog({
    kind: 'alert',
    title: title || titles[tone],
    message,
    tone,
    confirmLabel: 'Mengerti',
    cancelLabel: ''
  })
}

export const showAppConfirm = (
  message: string,
  options: {
    title?: string
    confirmLabel?: string
    cancelLabel?: string
    tone?: DialogTone
  } = {}
): Promise<boolean> =>
  enqueueDialog({
    kind: 'confirm',
    title: options.title || 'Konfirmasi',
    message,
    tone: options.tone || 'warning',
    confirmLabel: options.confirmLabel || 'Ya, lanjutkan',
    cancelLabel: options.cancelLabel || 'Batal'
  })

export const appDialogQueue = dialogQueue

export const resolveAppDialog = (confirmed: boolean): void => {
  const currentDialog = dialogQueue.value.shift()
  currentDialog?.resolve(confirmed)
}
