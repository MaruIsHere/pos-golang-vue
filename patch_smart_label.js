const fs = require('fs');
const file = 'frontend/src/components/PaymentModal.vue';
let code = fs.readFileSync(file, 'utf8');

const oldLabel = `<p class="text-[11px] text-slate-500 dark:text-slate-400">
                Pilih file, ambil foto, atau tekan <strong class="text-indigo-600 dark:text-indigo-400">Ctrl+V</strong> untuk Paste gambar (Maksimal 10 MB).
              </p>`;

const newLabel = `<p class="text-[11px] text-slate-500 dark:text-slate-400" v-html="smartUploadLabel"></p>`;

code = code.replace(oldLabel, newLabel);

const scriptInsert = `const smartUploadLabel = computed(() => {
  const ua = navigator.userAgent.toLowerCase();
  const isMobile = /android|webos|iphone|ipad|ipod|blackberry|windows phone/.test(ua);
  
  if (isMobile) {
    return 'Ambil foto struk dari kamera atau pilih gambar dari galeri HP (Maksimal 10 MB).';
  } else if (/macintosh|mac os x/.test(ua)) {
    return 'Pilih file atau tekan <strong class="text-indigo-600 dark:text-indigo-400">Cmd+V</strong> untuk langsung <em>Paste</em> gambar (Maks. 10 MB).';
  } else {
    return 'Pilih file atau tekan <strong class="text-indigo-600 dark:text-indigo-400">Ctrl+V</strong> untuk langsung <em>Paste</em> gambar (Maks. 10 MB).';
  }
});
`;

// Insert the computed property before onMounted
code = code.replace("onMounted(() => {", scriptInsert + "\nonMounted(() => {");

fs.writeFileSync(file, code);
