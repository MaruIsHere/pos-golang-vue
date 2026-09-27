import * as XLSX from 'xlsx';
import { jsPDF } from 'jspdf';
import autoTable from 'jspdf-autotable';
import type { DashboardStats } from '../types';

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

export const exportToExcel = (stats: DashboardStats, dateTitle: string = 'Keseluruhan'): void => {
  const wb = XLSX.utils.book_new();

  // 1. Sheet Ringkasan (Executive Summary)
  const storeName = stats.store_setting?.store_name || 'KASIR POS';
  const summaryData = [
    ['LAPORAN PENJUALAN KASIR POS'],
    ['Nama Toko:', storeName],
    ['Tanggal Cetak:', new Date().toLocaleString('id-ID')],
    ['Periode Laporan:', dateTitle],
    [''],
    ['--- RINGKASAN PERFORMA ---'],
    ['Total Omset Penjualan', stats.total_revenue],
    ['Total Transaksi', stats.total_orders],
    ['Total Item Terjual (Pcs)', stats.total_items_sold],
    ['']
  ];

  const wsSummary = XLSX.utils.aoa_to_sheet(summaryData);
  XLSX.utils.book_append_sheet(wb, wsSummary, 'Ringkasan Eksekutif');

  // 2. Sheet Barang Paling Laku (Top Sellers)
  const topRows = (stats.top_products || []).map((p, idx) => ({
    Peringkat: `#${idx + 1}`,
    'Nama Produk': p.product_name,
    Artist: p.artist || '-',
    'Tipe Produk': p.product_type || '-',
    'Harga (Rp)': p.price || 0,
    'Qty Terjual (Pcs)': p.total_qty,
    'Total Omset (Rp)': p.total_sales,
    'Sisa Stok': p.stock ?? '-'
  }));
  const wsTop = XLSX.utils.json_to_sheet(topRows);
  XLSX.utils.book_append_sheet(wb, wsTop, 'Barang Paling Laku');

  // 3. Sheet Barang Kurang Laku (Slow Moving)
  const leastRows = (stats.least_products || []).map((p, idx) => ({
    Peringkat: `#${idx + 1}`,
    'Nama Produk': p.product_name,
    Artist: p.artist || '-',
    'Tipe Produk': p.product_type || '-',
    'Harga (Rp)': p.price || 0,
    'Qty Terjual (Pcs)': p.total_qty,
    'Total Omset (Rp)': p.total_sales,
    'Sisa Stok Tersisa': p.stock ?? 0
  }));
  const wsLeast = XLSX.utils.json_to_sheet(wsLeastRows(leastRows));
  XLSX.utils.book_append_sheet(wb, wsLeast, 'Barang Kurang Laku');

  // 4. Sheet List Seluruh Barang Laku
  const allSoldRows = (stats.all_sold_products || []).map((p, idx) => ({
    No: idx + 1,
    'Nama Produk': p.product_name,
    Artist: p.artist || '-',
    'Tipe Produk': p.product_type || '-',
    'Harga (Rp)': p.price || 0,
    'Total Terjual (Pcs)': p.total_qty,
    'Total Revenue (Rp)': p.total_sales
  }));
  const wsAllSold = XLSX.utils.json_to_sheet(allSoldRows);
  XLSX.utils.book_append_sheet(wb, wsAllSold, 'List Seluruh Barang Laku');

  // 5. Sheet Breakdown Per Artist & Tipe
  const artistRows = (stats.sales_by_artist || []).map(a => ({
    'Nama Artist': a.name,
    'Total Terjual (Pcs)': a.total_qty,
    'Total Omset (Rp)': a.total_sales
  }));
  const wsArtist = XLSX.utils.json_to_sheet(artistRows);
  XLSX.utils.book_append_sheet(wb, wsArtist, 'Penjualan Per Artist');

  const typeRows = (stats.sales_by_type || []).map(t => ({
    'Tipe Produk': t.name,
    'Total Terjual (Pcs)': t.total_qty,
    'Total Omset (Rp)': t.total_sales
  }));
  const wsType = XLSX.utils.json_to_sheet(typeRows);
  XLSX.utils.book_append_sheet(wb, wsType, 'Penjualan Per Tipe');

  // Generate XLSX File
  const filename = `Laporan_Penjualan_${new Date().toISOString().slice(0, 10)}.xlsx`;
  XLSX.writeFile(wb, filename);
};

function wsLeastRows(rows: any[]) {
  return rows;
}

export const exportToPDF = (stats: DashboardStats, dateTitle: string = 'Keseluruhan'): void => {
  const doc = new jsPDF('p', 'mm', 'a4');
  const store = stats.store_setting;
  const storeName = store?.store_name || 'KASIR POS SYSTEM';
  const storeAddress = store?.address || 'Jl. Utama POS';
  const storePhone = store?.phone ? `Telp: ${store.phone}` : '';

  // Header Banner
  doc.setFillColor(30, 41, 59); // Slate-800
  doc.rect(0, 0, 210, 32, 'F');

  doc.setTextColor(255, 255, 255);
  doc.setFontSize(16);
  doc.setFont('helvetica', 'bold');
  doc.text(storeName.toUpperCase(), 14, 14);

  doc.setFontSize(8);
  doc.setFont('helvetica', 'normal');
  doc.text(`${storeAddress} ${storePhone}`, 14, 20);
  doc.text(`LAPORAN PERFORMA PENJUALAN PRODUK (${dateTitle.toUpperCase()})`, 14, 26);

  doc.setFontSize(8);
  doc.text(`Cetak: ${new Date().toLocaleString('id-ID')}`, 196, 26, { align: 'right' });

  // Key Metrics Cards
  let startY = 38;

  // Stat Card 1: Omset
  doc.setFillColor(240, 253, 244); // Light Green
  doc.setDrawColor(187, 247, 208);
  doc.roundedRect(14, startY, 58, 20, 3, 3, 'FD');
  doc.setTextColor(22, 101, 52);
  doc.setFontSize(8);
  doc.setFont('helvetica', 'bold');
  doc.text('TOTAL OMSET', 18, startY + 6);
  doc.setFontSize(11);
  doc.text(`Rp ${formatPrice(stats.total_revenue)}`, 18, startY + 14);

  // Stat Card 2: Transaksi
  doc.setFillColor(238, 242, 255); // Light Indigo
  doc.setDrawColor(199, 210, 254);
  doc.roundedRect(76, startY, 58, 20, 3, 3, 'FD');
  doc.setTextColor(55, 48, 163);
  doc.setFontSize(8);
  doc.setFont('helvetica', 'bold');
  doc.text('TOTAL TRANSAKSI', 80, startY + 6);
  doc.setFontSize(11);
  doc.text(`${stats.total_orders} Transaksi`, 80, startY + 14);

  // Stat Card 3: Item Terjual
  doc.setFillColor(254, 243, 199); // Light Amber
  doc.setDrawColor(253, 230, 138);
  doc.roundedRect(138, startY, 58, 20, 3, 3, 'FD');
  doc.setTextColor(146, 64, 14);
  doc.setFontSize(8);
  doc.setFont('helvetica', 'bold');
  doc.text('ITEM TERJUAL', 142, startY + 6);
  doc.setFontSize(11);
  doc.text(`${stats.total_items_sold} Pcs`, 142, startY + 14);

  startY += 26;

  // 1. Table 1: Barang Paling Laku (Top 10 Best Sellers)
  doc.setTextColor(15, 23, 42);
  doc.setFontSize(11);
  doc.setFont('helvetica', 'bold');
  doc.text('🔥 10 BARANG PALING LAKU (TOP SELLERS)', 14, startY);

  const topTableData = (stats.top_products || []).map((p, idx) => [
    `#${idx + 1}`,
    p.product_name,
    p.artist || 'Umum',
    p.product_type || 'Umum',
    `Rp ${formatPrice(p.price || 0)}`,
    `${p.total_qty} pcs`,
    `Rp ${formatPrice(p.total_sales)}`
  ]);

  autoTable(doc, {
    startY: startY + 3,
    head: [['Rank', 'Nama Produk', 'Artist', 'Tipe Produk', 'Harga', 'Terjual', 'Total Omset']],
    body: topTableData,
    theme: 'striped',
    headStyles: { fillColor: [79, 70, 229], fontSize: 8, fontStyle: 'bold' },
    bodyStyles: { fontSize: 8 },
    alternateRowStyles: { fillColor: [248, 250, 252] },
    margin: { left: 14, right: 14 }
  });

  // @ts-ignore
  startY = doc.lastAutoTable.finalY + 10;

  // 2. Table 2: Barang Kurang Laku (Slow Moving)
  doc.setFontSize(11);
  doc.setFont('helvetica', 'bold');
  doc.text('⚠️ 10 BARANG KURANG LAKU (SLOW MOVING / EVALUASI STOK)', 14, startY);

  const leastTableData = (stats.least_products || []).map((p, idx) => [
    `#${idx + 1}`,
    p.product_name,
    p.artist || 'Umum',
    p.product_type || 'Umum',
    `Rp ${formatPrice(p.price || 0)}`,
    `${p.total_qty} pcs`,
    `Sisa ${p.stock ?? 0} unit`
  ]);

  autoTable(doc, {
    startY: startY + 3,
    head: [['No', 'Nama Produk', 'Artist', 'Tipe Produk', 'Harga', 'Terjual', 'Stok Gudang']],
    body: leastTableData,
    theme: 'striped',
    headStyles: { fillColor: [225, 29, 72], fontSize: 8, fontStyle: 'bold' },
    bodyStyles: { fontSize: 8 },
    alternateRowStyles: { fillColor: [255, 241, 242] },
    margin: { left: 14, right: 14 }
  });

  // @ts-ignore
  startY = doc.lastAutoTable.finalY + 10;

  // Check page break if needed
  if (startY > 230) {
    doc.addPage();
    startY = 20;
  }

  // 3. Table 3: List Seluruh Barang Laku (Full Sales List)
  doc.setFontSize(11);
  doc.setFont('helvetica', 'bold');
  doc.text('📋 LIST SELURUH BARANG LAKU & PENJUALAN', 14, startY);

  const allSoldTableData = (stats.all_sold_products || []).map((p, idx) => [
    `${idx + 1}`,
    p.product_name,
    p.artist || 'Umum',
    p.product_type || 'Umum',
    `Rp ${formatPrice(p.price || 0)}`,
    `${p.total_qty} pcs`,
    `Rp ${formatPrice(p.total_sales)}`
  ]);

  autoTable(doc, {
    startY: startY + 3,
    head: [['No', 'Nama Produk', 'Artist', 'Tipe', 'Harga Satuan', 'Kuantitas', 'Subtotal Sales']],
    body: allSoldTableData,
    theme: 'grid',
    headStyles: { fillColor: [15, 23, 42], fontSize: 8, fontStyle: 'bold' },
    bodyStyles: { fontSize: 8 },
    margin: { left: 14, right: 14 }
  });

  // Add Page Numbers Footer
  const pageCount = (doc as any).internal.getNumberOfPages();
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i);
    doc.setFontSize(7);
    doc.setTextColor(148, 163, 184);
    doc.text(`Halaman ${i} dari ${pageCount} — Laporan Resmi ${storeName}`, 105, 290, { align: 'center' });
  }

  // Save PDF
  const filename = `Laporan_Penjualan_${new Date().toISOString().slice(0, 10)}.pdf`;
  doc.save(filename);
};
