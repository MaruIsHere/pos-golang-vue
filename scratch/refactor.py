import os
import re

files = [
    'frontend/src/views/OrdersView.vue',
    'frontend/src/views/ProductsView.vue',
    'frontend/src/views/CustomersView.vue',
    'frontend/src/views/InventoryView.vue'
]

def process_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Add missing imports for AppButton and AppInput
    if '<AppButton' not in content or '<AppInput' not in content:
        script_setup_idx = content.find('<script setup')
        if script_setup_idx != -1:
            end_script_tag = content.find('>', script_setup_idx) + 1
            imports = "\nimport AppButton from '@/components/ui/AppButton.vue';\nimport AppInput from '@/components/ui/AppInput.vue';\n"
            content = content[:end_script_tag] + imports + content[end_script_tag:]

    # 1. Replace `<button class="btn...">` with `<AppButton variant="...">`
    # We will do this carefully with regex.
    # Replace `<button class="btn btn-primary ..."`
    
    # We also need to change `</button>` to `</AppButton>` but only for those we changed. 
    # Because there might be generic `<button>` without `.btn` which we shouldn't change unless specified. 
    # But wait, instruction says "Replace raw <button class="btn..."> with <AppButton variant="...">" 
    # Let's just replace ALL `<button` that have `class="...btn..."` and then replace all `</button>` with `</AppButton>`. 
    # Wait, the safest is to replace all `<button` with `<AppButton` and `</button>` with `</AppButton>`, then fix the classes.
    content = re.sub(r'<button\b', '<AppButton', content)
    content = re.sub(r'</button>', '</AppButton>', content)

    # Now parse `<AppButton class="...">` to extract btn variants
    def btn_replacer(match):
        attrs = match.group(1)
        # check for btn-*
        variant = 'primary' # default
        if 'btn-primary' in attrs: variant = 'primary'
        elif 'btn-secondary' in attrs: variant = 'secondary'
        elif 'btn-success' in attrs: variant = 'success'
        elif 'btn-danger' in attrs: variant = 'danger'
        elif 'tab-btn' in attrs: variant = 'secondary' # generic mapping
        elif 'btn-icon' in attrs: variant = 'secondary'
        
        # remove btn, btn-primary, btn-secondary, btn-success, btn-danger, btn-icon, tab-btn
        new_attrs = re.sub(r'\b(btn|btn-primary|btn-secondary|btn-success|btn-danger|btn-icon|tab-btn|btn-submit|btn-edit|btn-delete|action-btns)\b', '', attrs)
        new_attrs = re.sub(r'class="\s+"', '', new_attrs)
        new_attrs = new_attrs.replace('class=" ', 'class="')
        
        return f'<AppButton variant="{variant}"{new_attrs}'

    content = re.sub(r'<AppButton([^>]*?class="[^"]*btn[^"]*"[^>]*?)', btn_replacer, content)

    # 2. Replace `<input class="form-control">` with `<AppInput>`
    # We need to change `<input` to `<AppInput` and self close if needed, but Vue allows normal `<AppInput>`
    content = re.sub(r'<input\b', '<AppInput', content)

    # 3. Replace legacy backgrounds/colors
    content = content.replace('bg-bg-card', 'bg-white dark:bg-slate-800')
    content = content.replace('bg-bg-primary', 'bg-slate-50 dark:bg-slate-900')
    content = content.replace('bg-card', 'bg-white dark:bg-slate-800')
    content = content.replace('bg-primary', 'bg-slate-50 dark:bg-slate-900')
    
    # 4. Replace .text-text-primary with text-slate-900 dark:text-slate-100
    content = content.replace('text-text-primary', 'text-slate-900 dark:text-slate-100')
    content = content.replace('text-primary', 'text-slate-900 dark:text-slate-100')
    content = content.replace('text-secondary', 'text-slate-500 dark:text-slate-400')
    content = content.replace('text-muted', 'text-slate-400 dark:text-slate-500')
    
    # 5. Replace .border-border-color with border-slate-200 dark:border-slate-700
    content = content.replace('border-border-color', 'border-slate-200 dark:border-slate-700')
    content = content.replace('border-color', 'border-slate-200 dark:border-slate-700')
    
    # 6. Replace .glass-panel
    content = content.replace('glass-panel', 'backdrop-blur-md bg-white/90 dark:bg-slate-900/90')

    # Remove all legacy classes from HTML elements
    legacy_classes = [
        'form-control', 'search-input', 'filter-select', 'products-page', 'orders-page',
        'customers-page', 'inventory-page', 'page-header', 'header-title', 'header-actions',
        'filter-bar', 'table-container', 'table-wrapper', 'data-table', 'movement-table', 
        'inventory-tabs', 'inventory-card', 'card-header', 'form-grid', 'full-width',
        'loading-box', 'empty-box', 'spinner', 'prod-thumb', 'prod-name', 'type-badge',
        'qty-badge', 'reason-tag', 'date-col', 'notes-col', 'in', 'out'
    ]

    # Map them to tailwind
    class_mappings = {
        'products-page': 'flex flex-col gap-5',
        'orders-page': 'flex flex-col gap-5',
        'customers-page': 'flex flex-col gap-5',
        'inventory-page': 'flex flex-col gap-5',
        'page-header': 'flex justify-between items-center p-5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl',
        'header-title': 'flex flex-col gap-1',
        'header-actions': 'flex items-center gap-3',
        'filter-bar': 'flex flex-wrap items-center gap-3 p-4 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl',
        'table-container': 'overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl',
        'table-wrapper': 'overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl',
        'data-table': 'w-full text-left text-sm',
        'movement-table': 'w-full text-left text-sm',
        'inventory-tabs': 'flex gap-2 overflow-x-auto shrink-0',
        'inventory-card': 'flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl',
        'card-header': 'border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100',
        'form-grid': 'grid grid-cols-1 sm:grid-cols-2 gap-4',
        'full-width': 'col-span-1 sm:col-span-2',
        'loading-box': 'p-12 text-center text-slate-500',
        'empty-box': 'p-12 text-center text-slate-500',
        'spinner': 'w-8 h-8 border-4 border-slate-200 dark:border-slate-700 border-t-indigo-600 rounded-full animate-spin mx-auto mb-2',
        'prod-thumb': 'w-10 h-10 rounded-lg object-cover border border-slate-200 dark:border-slate-700',
        'prod-name': 'font-semibold text-slate-900 dark:text-slate-100',
        'type-badge': 'px-2.5 py-0.5 rounded-full text-xs font-bold',
        'qty-badge': 'font-bold text-sm',
        'reason-tag': 'px-2 py-0.5 rounded text-xs bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400',
        'date-col': 'text-xs text-slate-500',
        'notes-col': 'text-xs text-slate-500',
        'form-control': 'w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100',
        'search-input': 'flex-auto min-w-[200px]',
        'filter-select': 'flex-none min-w-[140px]',
    }

    # Add text mappings for h2, p, th, td
    content = content.replace('<h2>', '<h2 class="text-xl font-bold text-slate-900 dark:text-slate-100">')
    content = content.replace('<h3>', '<h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">')
    content = content.replace('<p>', '<p class="text-sm text-slate-500 dark:text-slate-400">')
    content = content.replace('<th>', '<th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">')
    content = content.replace('<td>', '<td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">')

    # Also clean classes
    for old_class, new_class in class_mappings.items():
        # Only replace word bounded matches in class attributes
        # A bit risky, but let's just do a simple replacement
        # We can just replace the strings directly if they are inside class=""
        def class_replacer(match):
            cls_str = match.group(1)
            # Replace exact words
            cls_str = re.sub(rf'\b{old_class}\b', new_class, cls_str)
            return f'class="{cls_str}"'
            
        content = re.sub(r'class="([^"]+)"', class_replacer, content)

    # 7. Remove all <style scoped> blocks
    content = re.sub(r'<style scoped>.*?</style>', '', content, flags=re.DOTALL)
    # Also remove `<style>` blocks just in case
    content = re.sub(r'<style>.*?</style>', '', content, flags=re.DOTALL)

    # Save
    with open(filepath, 'w') as f:
        f.write(content)
    print(f"Processed {filepath}")

for f in files:
    process_file(f)
