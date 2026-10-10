const fs = require('fs');
const file = 'frontend/src/components/PaymentModal.vue';
let code = fs.readFileSync(file, 'utf8');

const oldButton = `<button 
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
              </button>`;

const newButton = `<input 
                  id="proof-upload-input"
                  ref="proofInputRef"
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  @change="onProofFileSelected" 
                  :disabled="isUploadingProof"
                />
              <button 
                type="button"
                @click="triggerProofUpload"
                class="w-full py-2.5 px-4 border-2 border-dashed border-indigo-300 dark:border-indigo-700 hover:border-indigo-500 dark:hover:border-indigo-500 rounded-xl bg-white dark:bg-slate-900 flex items-center justify-center gap-2 text-xs font-bold text-indigo-700 dark:text-indigo-300 transition-colors shadow-xs cursor-pointer"
                :class="{ 'opacity-50 cursor-not-allowed': isUploadingProof }"
                :disabled="isUploadingProof"
              >
                <PhotoIcon class="w-4 h-4 text-indigo-500" />
                <span>{{ isUploadingProof ? 'Mengunggah Foto...' : (paymentProofPreview ? 'Ganti Foto Bukti Bayar' : 'Upload / Ambil Foto Bukti Pembayaran') }}</span>
              </button>`;

code = code.replace(oldButton, newButton);
fs.writeFileSync(file, code);
