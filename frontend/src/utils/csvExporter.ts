import type { DashboardStats } from '../types';

export const exportToCSV = (stats: DashboardStats, dateTitle: string = 'Keseluruhan'): void => {
  const allSold = stats.all_sold_products || [];
  
  const headers = ['No', 'Nama Produk', 'Merk/Artist', 'Tipe Produk', 'Harga Satuan', 'Kuantitas Terjual', 'Total Omset'];
  const csvRows = [headers.join(',')];

  // Prevent CSV Injection by prefixing potentially dangerous cells with a single quote
  const sanitize = (cell: string | number) => {
    let str = String(cell);
    if (/^[=+\-@]/.test(str)) {
      str = "'" + str;
    }
    return `"${str.replace(/"/g, '""')}"`;
  };
  
  allSold.forEach((p, idx) => {
    const row = [
      idx + 1,
      sanitize(p.product_name || ''),
      sanitize(p.artist || 'Umum'),
      sanitize(p.product_type || 'Umum'),
      p.price || 0,
      p.total_qty || 0,
      p.total_sales || 0
    ];
    csvRows.push(row.join(','));
  });

  // Add UTF-8 BOM for Excel compatibility
  const csvString = '\uFEFF' + csvRows.join('\n');
  const blob = new Blob([csvString], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `Laporan_Penjualan_${new Date().toISOString().slice(0, 10)}.csv`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};
