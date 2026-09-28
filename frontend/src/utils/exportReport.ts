import * as XLSX from 'xlsx-js-style/dist/xlsx.bundle.js';
import { jsPDF } from 'jspdf';
import autoTable from 'jspdf-autotable';
import type { DashboardStats } from '../types';

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const appendTableSheet = (
  workbook: XLSX.WorkBook,
  name: string,
  rows: Record<string, string | number>[],
  headers: string[],
  widths: number[],
  currencyHeaders: string[] = []
): void => {
  const worksheet = XLSX.utils.json_to_sheet(rows, { header: headers });
  const range = XLSX.utils.decode_range(worksheet['!ref'] || 'A1');
  worksheet['!cols'] = widths.map(wch => ({ wch }));
  worksheet['!autofilter'] = { ref: XLSX.utils.encode_range(range) };

  for (let row = range.s.r; row <= range.e.r; row++) {
    for (let column = range.s.c; column <= range.e.c; column++) {
      const address = XLSX.utils.encode_cell({ r: row, c: column });
      const cell = worksheet[address];
      if (!cell) continue;

      if (row === 0) {
        cell.s = {
          fill: { patternType: 'solid', fgColor: { rgb: '176B5B' } },
          font: { bold: true, color: { rgb: 'FFFFFF' } },
          alignment: { vertical: 'center', wrapText: true }
        };
      } else {
        cell.s = {
          fill: { patternType: 'solid', fgColor: { rgb: row % 2 === 0 ? 'F1F7F5' : 'FFFFFF' } },
          font: { color: { rgb: '24332F' } },
          alignment: { vertical: 'center' }
        };
        if (currencyHeaders.includes(headers[column]) && typeof cell.v === 'number') {
          cell.z = '"Rp" #,##0';
          cell.s.alignment = { vertical: 'center', horizontal: 'right' };
        } else if (typeof cell.v === 'number') {
          cell.s.alignment = { vertical: 'center', horizontal: 'right' };
        }
      }
    }
  }

  worksheet['!rows'] = [{ hpt: 30 }];
  XLSX.utils.book_append_sheet(workbook, worksheet, name);
};

const addSalesCharts = async (workbook: XLSX.WorkBook, stats: DashboardStats): Promise<void> => {
  const [{ Chart, registerables }, ExcelJS] = await Promise.all([
    import('chart.js'),
    import('exceljs')
  ]);
  Chart.register(...registerables);

  const excelWorkbook = new ExcelJS.Workbook();
  const xlsxBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });
  await excelWorkbook.xlsx.load(xlsxBuffer as ArrayBuffer);

  const sheet = excelWorkbook.addWorksheet('Grafik Penjualan');
  sheet.columns = Array.from({ length: 16 }, () => ({ width: 12 }));
  sheet.views = [{ showGridLines: false }];
  sheet.mergeCells('A1:P1');
  sheet.getCell('A1').value = 'VISUALISASI PENJUALAN';
  sheet.getCell('A1').font = { bold: true, size: 16, color: { argb: 'FFFFFFFF' } };
  sheet.getCell('A1').fill = { type: 'pattern', pattern: 'solid', fgColor: { argb: 'FF174C43' } };
  sheet.getRow(1).height = 34;
  for (let row = 3; row <= 21; row++) sheet.getRow(row).height = 22;

  const renderChart = (type: 'bar' | 'doughnut', labels: string[], values: number[], title: string): string => {
    const canvas = document.createElement('canvas');
    canvas.width = 1000;
    canvas.height = 560;
    const chart = new Chart(canvas, {
      type,
      data: {
        labels,
        datasets: [{
          label: title,
          data: values,
          backgroundColor: type === 'bar'
            ? '#176B5B'
            : ['#176B5B', '#E58A4E', '#3D8CBB', '#D6B84C', '#C45E64', '#638B73'],
          borderColor: '#FFFFFF',
          borderWidth: 2
        }]
      },
      options: {
        responsive: false,
        animation: false,
        ...(type === 'bar' ? { indexAxis: 'y' as const } : {}),
        plugins: {
          title: { display: true, text: title },
          legend: { display: type === 'doughnut', position: 'bottom' as const }
        }
      }
    } as import('chart.js').ChartConfiguration);
    const image = chart.toBase64Image();
    chart.destroy();
    return image.split(',')[1];
  };

  const topProducts = (stats.top_products || []).filter(product => product.total_qty > 0).slice(0, 8).reverse();
  if (topProducts.length) {
    const imageId = excelWorkbook.addImage({
      base64: renderChart('bar', topProducts.map(product => product.product_name), topProducts.map(product => product.total_qty), 'Produk Terlaris (Jumlah Terjual)'),
      extension: 'png'
    });
    sheet.addImage(imageId, 'A3:H21');
  } else {
    sheet.getCell('A3').value = 'Belum ada data produk terjual.';
  }

  const productTypes = (stats.sales_by_type || []).filter(item => item.total_sales > 0);
  if (productTypes.length) {
    const imageId = excelWorkbook.addImage({
      base64: renderChart('doughnut', productTypes.map(item => item.name), productTypes.map(item => item.total_sales), 'Omzet per Tipe Produk'),
      extension: 'png'
    });
    sheet.addImage(imageId, 'I3:P21');
  } else {
    sheet.getCell('I3').value = 'Belum ada data omzet per tipe produk.';
  }

  const output = await excelWorkbook.xlsx.writeBuffer();
  const blob = new Blob([output as BlobPart], {
    type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
  });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `Laporan_Penjualan_${new Date().toISOString().slice(0, 10)}.xlsx`;
  link.click();
  window.setTimeout(() => URL.revokeObjectURL(url), 1000);
};

export const exportToExcel = async (stats: DashboardStats, dateTitle: string = 'Keseluruhan'): Promise<void> => {
  const wb = XLSX.utils.book_new();

  // 1. Sheet Ringkasan (Executive Summary)
  const storeName = stats.store_setting?.store_name || 'KASIR POS';
  const summaryData = [
    ['LAPORAN PENJUALAN'],
    ['Nama Toko:', storeName],
    ['Tanggal Cetak:', new Date().toLocaleString('id-ID')],
    ['Periode Laporan:', dateTitle],
    [''],
    ['RINGKASAN PERFORMA'],
    ['Total Omset Penjualan', stats.total_revenue],
    ['Total Transaksi', stats.total_orders],
    ['Total Item Terjual (Pcs)', stats.total_items_sold],
    ['']
  ];

  const wsSummary = XLSX.utils.aoa_to_sheet(summaryData);
  wsSummary['!merges'] = [
    { s: { r: 0, c: 0 }, e: { r: 0, c: 1 } },
    { s: { r: 5, c: 0 }, e: { r: 5, c: 1 } }
  ];
  wsSummary['!cols'] = [{ wch: 34 }, { wch: 38 }];
  wsSummary['!rows'] = [{ hpt: 34 }, { hpt: 23 }, { hpt: 23 }, { hpt: 23 }, { hpt: 10 }, { hpt: 25 }];
  wsSummary['A1'].s = {
    fill: { patternType: 'solid', fgColor: { rgb: '174C43' } },
    font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 16 },
    alignment: { vertical: 'center', horizontal: 'left' }
  };
  wsSummary['A6'].s = {
    fill: { patternType: 'solid', fgColor: { rgb: '176B5B' } },
    font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 11 },
    alignment: { vertical: 'center' }
  };
  for (let row = 1; row <= 3; row++) {
    wsSummary[`A${row + 1}`].s = {
      fill: { patternType: 'solid', fgColor: { rgb: 'EEF5F2' } },
      font: { bold: true, color: { rgb: '334A43' } },
      alignment: { vertical: 'center' }
    };
    wsSummary[`B${row + 1}`].s = {
      font: { color: { rgb: '24332F' } },
      alignment: { vertical: 'center' }
    };
  }
  for (let row = 7; row <= 9; row++) {
    wsSummary[`A${row}`].s = {
      fill: { patternType: 'solid', fgColor: { rgb: 'EEF5F2' } },
      font: { bold: true, color: { rgb: '334A43' } },
      alignment: { vertical: 'center' }
    };
    wsSummary[`B${row}`].s = {
      fill: { patternType: 'solid', fgColor: { rgb: 'F8FBF9' } },
      font: { bold: true, color: { rgb: '176B5B' }, sz: 12 },
      alignment: { vertical: 'center', horizontal: 'right' },
      ...(row === 7 ? { numFmt: '"Rp" #,##0' } : {})
    };
  }
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
  appendTableSheet(wb, 'Barang Paling Laku', topRows, [
    'Peringkat', 'Nama Produk', 'Artist', 'Tipe Produk', 'Harga (Rp)', 'Qty Terjual (Pcs)', 'Total Omset (Rp)', 'Sisa Stok'
  ], [12, 34, 22, 20, 17, 20, 20, 14], ['Harga (Rp)', 'Total Omset (Rp)']);

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
  appendTableSheet(wb, 'Barang Kurang Laku', leastRows, [
    'Peringkat', 'Nama Produk', 'Artist', 'Tipe Produk', 'Harga (Rp)', 'Qty Terjual (Pcs)', 'Total Omset (Rp)', 'Sisa Stok Tersisa'
  ], [12, 34, 22, 20, 17, 20, 20, 20], ['Harga (Rp)', 'Total Omset (Rp)']);

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
  appendTableSheet(wb, 'List Seluruh Barang Laku', allSoldRows, [
    'No', 'Nama Produk', 'Artist', 'Tipe Produk', 'Harga (Rp)', 'Total Terjual (Pcs)', 'Total Revenue (Rp)'
  ], [10, 36, 22, 20, 17, 22, 22], ['Harga (Rp)', 'Total Revenue (Rp)']);

  // 5. Sheet Breakdown Per Artist & Tipe
  const artistRows = (stats.sales_by_artist || []).map(a => ({
    'Nama Artist': a.name,
    'Total Terjual (Pcs)': a.total_qty,
    'Total Omset (Rp)': a.total_sales
  }));
  appendTableSheet(wb, 'Penjualan Per Artist', artistRows, [
    'Nama Artist', 'Total Terjual (Pcs)', 'Total Omset (Rp)'
  ], [34, 23, 23], ['Total Omset (Rp)']);

  const typeRows = (stats.sales_by_type || []).map(t => ({
    'Tipe Produk': t.name,
    'Total Terjual (Pcs)': t.total_qty,
    'Total Omset (Rp)': t.total_sales
  }));
  appendTableSheet(wb, 'Penjualan Per Tipe', typeRows, [
    'Tipe Produk', 'Total Terjual (Pcs)', 'Total Omset (Rp)'
  ], [34, 23, 23], ['Total Omset (Rp)']);

  await addSalesCharts(wb, stats);
};


const addPdfSalesCharts = (doc: jsPDF, stats: DashboardStats, dateTitle: string): void => {
  doc.addPage();
  doc.setFillColor(23, 76, 67);
  doc.rect(0, 0, 210, 32, 'F');
  doc.setTextColor(255, 255, 255);
  doc.setFont('helvetica', 'bold');
  doc.setFontSize(15);
  doc.text('VISUALISASI PENJUALAN', 14, 14);
  doc.setFont('helvetica', 'normal');
  doc.setFontSize(8);
  doc.text(`Periode: ${dateTitle}`, 14, 23);

  const drawBars = (
    title: string,
    items: { name: string; value: number }[],
    startY: number,
    color: [number, number, number]
  ): void => {
    doc.setTextColor(30, 41, 59);
    doc.setFont('helvetica', 'bold');
    doc.setFontSize(11);
    doc.text(title, 14, startY);

    if (!items.length) {
      doc.setFont('helvetica', 'normal');
      doc.setFontSize(9);
      doc.setTextColor(100, 116, 139);
      doc.text('Belum ada data penjualan.', 14, startY + 12);
      return;
    }

    const maxValue = Math.max(...items.map(item => item.value), 1);
    items.forEach((item, index) => {
      const rowY = startY + 9 + index * 14;
      const labelLines = doc.splitTextToSize(item.name, 54).slice(0, 2);
      doc.setFont('helvetica', 'normal');
      doc.setFontSize(8);
      doc.setTextColor(51, 65, 85);
      doc.text(labelLines, 14, rowY + 4);

      doc.setFillColor(241, 245, 249);
      doc.roundedRect(73, rowY, 82, 6, 1, 1, 'F');
      doc.setFillColor(...color);
      doc.roundedRect(73, rowY, Math.max((item.value / maxValue) * 82, 1), 6, 1, 1, 'F');

      doc.setTextColor(30, 41, 59);
      doc.setFontSize(7);
      doc.text(`Rp ${formatPrice(item.value)}`, 196, rowY + 4, { align: 'right', maxWidth: 38 });
    });
  };

  const topProducts = [...(stats.top_products || [])]
    .filter(product => product.total_sales > 0)
    .sort((left, right) => right.total_sales - left.total_sales)
    .slice(0, 5)
    .map(product => ({ name: product.product_name, value: product.total_sales }));
  const productTypes = [...(stats.sales_by_type || [])]
    .filter(type => type.total_sales > 0)
    .sort((left, right) => right.total_sales - left.total_sales)
    .slice(0, 5)
    .map(type => ({ name: type.name, value: type.total_sales }));

  drawBars('5 Produk dengan Omzet Tertinggi', topProducts, 44, [23, 107, 91]);
  drawBars('Omzet per Tipe Produk (5 Teratas)', productTypes, 137, [61, 140, 187]);
};

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

  addPdfSalesCharts(doc, stats, dateTitle);

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
