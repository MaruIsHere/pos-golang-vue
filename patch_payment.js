const fs = require('fs');
const file = 'frontend/src/components/PaymentModal.vue';
let code = fs.readFileSync(file, 'utf8');

// Replace the label with a button
const oldInputArea = `            <div class="flex-1 w-full flex flex-col gap-1.5">
              <label 
                for="proof-upload-input"
                class="w-full py-2.5 px-4 border-2 border-dashed border-indigo-300 dark:border-indigo-700 hover:border-indigo-500 dark:hover:border-indigo-500 rounded-xl bg-white dark:bg-slate-900 flex items-center justify-center gap-2 text-xs font-bold text-indigo-700 dark:text-indigo-300 transition-colors shadow-xs cursor-pointer"
                :class="{ 'opacity-50 cursor-not-allowed': isUploadingProof }"
              >
                <input 
                  id="proof-upload-input"
                  ref="proofInputRef"
                  type="file" 
                  accept="image/*" 
                  class="sr-only" 
                  @change="onProofFileSelected" 
                  :disabled="isUploadingProof"
                />
                <PhotoIcon class="w-4 h-4 text-indigo-500" />
                <span>{{ isUploadingProof ? 'Mengunggah Foto...' : (paymentProofPreview ? 'Ganti Foto Bukti Bayar' : 'Upload / Ambil Foto Bukti Pembayaran') }}</span>
              </label>
              
              <p class="text-[11px] text-slate-500 dark:text-slate-400">
                Pilih atau ambil foto bukti transfer / QRIS (Maksimal 10 MB).
              </p>
            </div>`;

const newInputArea = `            <div class="flex-1 w-full flex flex-col gap-1.5">
              <button 
                type="button"
                @click="triggerProofUpload"
                class="w-full py-2.5 px-4 border-2 border-dashed border-indigo-300 dark:border-indigo-700 hover:border-indigo-500 dark:hover:border-indigo-500 rounded-xl bg-white dark:bg-slate-900 flex items-center justify-center gap-2 text-xs font-bold text-indigo-700 dark:text-indigo-300 transition-colors shadow-xs cursor-pointer"
                :class="{ 'opacity-50 cursor-not-allowed': isUploadingProof }"
                :disabled="isUploadingProof"
              >
                <input 
                  id="proof-upload-input"
                  ref="proofInputRef"
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  @change="onProofFileSelected" 
                  :disabled="isUploadingProof"
                />
                <PhotoIcon class="w-4 h-4 text-indigo-500" />
                <span>{{ isUploadingProof ? 'Mengunggah Foto...' : (paymentProofPreview ? 'Ganti Foto Bukti Bayar' : 'Upload / Ambil Foto Bukti Pembayaran') }}</span>
              </button>
              
              <p class="text-[11px] text-slate-500 dark:text-slate-400">
                Pilih file, ambil foto, atau tekan <strong class="text-indigo-600 dark:text-indigo-400">Ctrl+V</strong> untuk Paste gambar (Maksimal 10 MB).
              </p>
            </div>`;

code = code.replace(oldInputArea, newInputArea);

// Extract the upload logic into a reusable function
const oldUploadFunction = `const onProofFileSelected = async (e: Event) => {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  paymentProofPreview.value = URL.createObjectURL(file);
  isUploadingProof.value = true;
  try {
    const formData = new FormData();
    formData.append('proof', file);
    const { data } = await api.post<{ payment_proof: string, detected_amount?: number }>('/orders/payment-proof', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    paymentProofUrl.value = data.payment_proof;
    
    // Auto-fill paid amount from OCR detection!
    if (data.detected_amount && data.detected_amount > 0) {
      paidAmount.value = data.detected_amount;
      console.log('OCR detected amount:', data.detected_amount);
    }
  } catch (err: any) {
    console.error('Error uploading payment proof:', err);
  } finally {
    isUploadingProof.value = false;
  }
};`;

const newUploadFunction = `const processProofFile = async (file: File) => {
  paymentProofPreview.value = URL.createObjectURL(file);
  isUploadingProof.value = true;
  try {
    const formData = new FormData();
    formData.append('proof', file);
    const { data } = await api.post<{ payment_proof: string, detected_amount?: number }>('/orders/payment-proof', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    paymentProofUrl.value = data.payment_proof;
    
    // Auto-fill paid amount from OCR detection!
    if (data.detected_amount && data.detected_amount > 0) {
      paidAmount.value = data.detected_amount;
      console.log('OCR detected amount:', data.detected_amount);
    }
  } catch (err: any) {
    console.error('Error uploading payment proof:', err);
  } finally {
    isUploadingProof.value = false;
  }
};

const onProofFileSelected = async (e: Event) => {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  await processProofFile(file);
};

// Clipboard Paste Support for Desktop/PC
const handlePaste = (e: ClipboardEvent) => {
  if (paymentMethod.value !== 'qris' && paymentMethod.value !== 'transfer') return;
  
  const items = e.clipboardData?.items;
  if (!items) return;
  
  for (let i = 0; i < items.length; i++) {
    if (items[i].type.indexOf('image') !== -1) {
      e.preventDefault();
      const blob = items[i].getAsFile();
      if (blob) {
        processProofFile(blob);
      }
      break;
    }
  }
};
`;

code = code.replace(oldUploadFunction, newUploadFunction);

// Add event listeners in onMounted
const oldOnMounted = `onMounted(() => {
  fetchCustomers();
  fetchLatestSettings();
});`;

const newOnMounted = `import { onBeforeUnmount } from 'vue';

onMounted(() => {
  fetchCustomers();
  fetchLatestSettings();
  window.addEventListener('paste', handlePaste);
});

onBeforeUnmount(() => {
  window.removeEventListener('paste', handlePaste);
});`;

// Replace onMounted, but be careful with imports
code = code.replace("import { ref, computed, onMounted } from 'vue';", "import { ref, computed, onMounted, onBeforeUnmount } from 'vue';");
code = code.replace(oldOnMounted, `onMounted(() => {
  fetchCustomers();
  fetchLatestSettings();
  window.addEventListener('paste', handlePaste);
});

onBeforeUnmount(() => {
  window.removeEventListener('paste', handlePaste);
});`);

fs.writeFileSync(file, code);
