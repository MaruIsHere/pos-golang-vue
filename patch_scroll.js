const fs = require('fs');
const file = 'frontend/src/components/PaymentModal.vue';
let code = fs.readFileSync(file, 'utf8');

// Add ids to sections
code = code.replace('<div class="p-4 sm:p-5 flex flex-col gap-4 overflow-y-auto flex-1 min-h-0">', '<div id="payment-scroll-container" class="p-4 sm:p-5 flex flex-col gap-4 overflow-y-auto flex-1 min-h-0 scroll-smooth">');

code = code.replace('<div v-if="paymentMethod === \'cash\'" class="bg-slate-50', '<div id="payment-section-cash" v-if="paymentMethod === \'cash\'" class="bg-slate-50');
code = code.replace('<div v-else-if="paymentMethod === \'qris\'" class="bg-slate-50', '<div id="payment-section-qris" v-else-if="paymentMethod === \'qris\'" class="bg-slate-50');
code = code.replace('<div v-else-if="paymentMethod === \'transfer\'" class="bg-slate-50', '<div id="payment-section-transfer" v-else-if="paymentMethod === \'transfer\'" class="bg-slate-50');
code = code.replace('<div v-else-if="paymentMethod === \'debit\'" class="bg-slate-50', '<div id="payment-section-debit" v-else-if="paymentMethod === \'debit\'" class="bg-slate-50');

// Add nextTick to vue imports
code = code.replace(/import { ref, computed, onMounted, onBeforeUnmount } from 'vue';/, "import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue';");

// Update selectMethod
const oldSelectMethod = `const selectMethod = (method: string): void => {
  paymentMethod.value = method;
  if (method !== 'cash') {
    paidAmount.value = finalGrandTotal.value;
  }
};`;

const newSelectMethod = `const selectMethod = async (method: string): Promise<void> => {
  paymentMethod.value = method;
  if (method !== 'cash') {
    paidAmount.value = finalGrandTotal.value;
  }
  
  await nextTick();
  
  // Auto-scroll logic to improve mobile UX
  const el = document.getElementById(\`payment-section-\${method}\`);
  const container = document.getElementById('payment-scroll-container');
  if (el && container) {
    // Scroll the container so the selected section is near the top
    container.scrollTo({
      top: el.offsetTop - 16,
      behavior: 'smooth'
    });
  }
};`;

code = code.replace(oldSelectMethod, newSelectMethod);

fs.writeFileSync(file, code);
