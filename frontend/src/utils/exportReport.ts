import * as XLSX from 'xlsx-js-style/dist/xlsx.bundle.js';
import { jsPDF } from 'jspdf';
import autoTable from 'jspdf-autotable';
import type { DashboardStats } from '../types';

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);
const formatQuantity = (val: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(val || 0);
const unitLabel = (unit?: string): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';

const appendTableSheet = (
  workbook: any,
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

const addSalesCharts = async (workbook: any, stats: DashboardStats): Promise<void> => {
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

  const store = stats.store_setting;
  const storeName = store?.store_name || 'KASIR POS SYSTEM';
  const storeAddress = store?.address || 'Jl. Utama POS';
  const storePhone = store?.phone ? `Telp: ${store.phone}` : '';
  const printedDate = new Date().toLocaleString('id-ID');

  // ==========================================
  // SHEET 1: Ringkasan & Top Items (Main Report)
  // ==========================================
  const sheet1Data: any[][] = [
    [storeName.toUpperCase()],
    [`${storeAddress} ${storePhone}`],
    [`LAPORAN PERFORMA PENJUALAN PRODUK (${dateTitle.toUpperCase()}) — Cetak: ${printedDate}`],
    [''], // Row 4 spacer

    // Stat Cards (Row 5 & 6)
    ['TOTAL OMSET', '', 'TOTAL TRANSAKSI', '', 'ITEM TERJUAL', '', ''],
    [stats.total_revenue || 0, '', `${stats.total_orders || 0} Transaksi`, '', `${formatQuantity(stats.total_items_sold || 0)} Item`, '', ''],
    [''], // Row 7 spacer

    // Table 1: Top Sellers (Row 8 Header, Row 9 Cols)
    ['🔥 10 BARANG PALING LAKU (TOP SELLERS)'],
    ['Rank', 'Nama Produk', 'Artist / Merk', 'Tipe Produk', 'Harga Satuan', 'Terjual', 'Total Omset']
  ];

  // Populate Top 10 Products
  const topList = (stats.top_products || []).slice(0, 10);
  topList.forEach((p, idx) => {
    sheet1Data.push([
      `#${idx + 1}`,
      p.product_name,
      p.artist || 'Umum',
      p.product_type || 'Umum',
      p.price || 0,
      `${formatQuantity(p.total_qty)} ${unitLabel(p.unit)}`,
      p.total_sales || 0
    ]);
  });

  // Table 2: Slow Moving
  sheet1Data.push(['']);
  const slowHeaderRowIdx = sheet1Data.length;
  sheet1Data.push(['⚠️ 10 BARANG KURANG LAKU (SLOW MOVING / EVALUASI STOK)']);
  const slowColsRowIdx = sheet1Data.length;
  sheet1Data.push(['No', 'Nama Produk', 'Artist / Merk', 'Tipe Produk', 'Harga Satuan', 'Terjual', 'Stok Gudang']);

  const leastList = (stats.least_products || []).slice(0, 10);
  leastList.forEach((p, idx) => {
    sheet1Data.push([
      `#${idx + 1}`,
      p.product_name,
      p.artist || 'Umum',
      p.product_type || 'Umum',
      p.price || 0,
      `${formatQuantity(p.total_qty)} ${unitLabel(p.unit)}`,
      `Sisa ${formatQuantity(p.stock ?? 0)} ${unitLabel(p.unit)}`
    ]);
  });

  const ws1 = XLSX.utils.aoa_to_sheet(sheet1Data);

  // Set Merges for Sheet 1
  ws1['!merges'] = [
    { s: { r: 0, c: 0 }, e: { r: 0, c: 6 } }, // Header Title
    { s: { r: 1, c: 0 }, e: { r: 1, c: 6 } }, // Header Address
    { s: { r: 2, c: 0 }, e: { r: 2, c: 6 } }, // Header Date
    // Stat cards merges
    { s: { r: 4, c: 0 }, e: { r: 4, c: 1 } },
    { s: { r: 5, c: 0 }, e: { r: 5, c: 1 } },
    { s: { r: 4, c: 2 }, e: { r: 4, c: 3 } },
    { s: { r: 5, c: 2 }, e: { r: 5, c: 3 } },
    { s: { r: 4, c: 4 }, e: { r: 4, c: 6 } },
    { s: { r: 5, c: 4 }, e: { r: 5, c: 6 } },
    // Table headers
    { s: { r: 7, c: 0 }, e: { r: 7, c: 6 } },
    { s: { r: slowHeaderRowIdx, c: 0 }, e: { r: slowHeaderRowIdx, c: 6 } }
  ];

  ws1['!cols'] = [
    { wch: 10 }, // Rank / No
    { wch: 34 }, // Nama Produk
    { wch: 20 }, // Artist / Merk
    { wch: 20 }, // Tipe Produk
    { wch: 18 }, // Harga Satuan
    { wch: 20 }, // Qty Terjual
    { wch: 22 }  // Total Omset / Stok
  ];

  // Apply Styling for Sheet 1
  if (ws1['A1']) ws1['A1'].s = { fill: { fgColor: { rgb: '1E293B' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 14 }, alignment: { vertical: 'center' } };
  if (ws1['A2']) ws1['A2'].s = { fill: { fgColor: { rgb: '1E293B' } }, font: { color: { rgb: 'CBD5E1' }, sz: 9 }, alignment: { vertical: 'center' } };
  if (ws1['A3']) ws1['A3'].s = { fill: { fgColor: { rgb: '1E293B' } }, font: { color: { rgb: '94A3B8' }, sz: 9, italic: true }, alignment: { vertical: 'center' } };

  // Stat Card 1 (Omset)
  if (ws1['A5']) ws1['A5'].s = { fill: { fgColor: { rgb: 'F0FDF4' } }, font: { bold: true, color: { rgb: '166534' }, sz: 9 }, alignment: { vertical: 'center', horizontal: 'center' } };
  if (ws1['A6']) {
    ws1['A6'].s = { fill: { fgColor: { rgb: 'F0FDF4' } }, font: { bold: true, color: { rgb: '166534' }, sz: 13 }, alignment: { vertical: 'center', horizontal: 'center' } };
    ws1['A6'].z = '"Rp" #,##0';
  }

  // Stat Card 2 (Transaksi)
  if (ws1['C5']) ws1['C5'].s = { fill: { fgColor: { rgb: 'EEF2FF' } }, font: { bold: true, color: { rgb: '3730A3' }, sz: 9 }, alignment: { vertical: 'center', horizontal: 'center' } };
  if (ws1['C6']) ws1['C6'].s = { fill: { fgColor: { rgb: 'EEF2FF' } }, font: { bold: true, color: { rgb: '3730A3' }, sz: 13 }, alignment: { vertical: 'center', horizontal: 'center' } };

  // Stat Card 3 (Item Terjual)
  if (ws1['E5']) ws1['E5'].s = { fill: { fgColor: { rgb: 'FEF3C7' } }, font: { bold: true, color: { rgb: '92400E' }, sz: 9 }, alignment: { vertical: 'center', horizontal: 'center' } };
  if (ws1['E6']) ws1['E6'].s = { fill: { fgColor: { rgb: 'FEF3C7' } }, font: { bold: true, color: { rgb: '92400E' }, sz: 13 }, alignment: { vertical: 'center', horizontal: 'center' } };

  // Section 1 Header (Top Sellers)
  if (ws1['A8']) ws1['A8'].s = { fill: { fgColor: { rgb: '4F46E5' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 11 }, alignment: { vertical: 'center' } };
  for (let c = 0; c < 7; c++) {
    const colName = String.fromCharCode(65 + c);
    if (ws1[`${colName}9`]) {
      ws1[`${colName}9`].s = { fill: { fgColor: { rgb: '3730A3' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 9 }, alignment: { vertical: 'center' } };
    }
  }

  // Section 1 Data rows
  topList.forEach((_, idx) => {
    const r = 10 + idx;
    const bg = idx % 2 === 0 ? 'F8FAFC' : 'FFFFFF';
    for (let c = 0; c < 7; c++) {
      const cellRef = `${String.fromCharCode(65 + c)}${r}`;
      if (ws1[cellRef]) {
        ws1[cellRef].s = { fill: { fgColor: { rgb: bg } }, font: { sz: 9 }, alignment: { vertical: 'center' } };
        if (c === 4 || c === 6) {
          ws1[cellRef].z = '"Rp" #,##0';
          ws1[cellRef].s.alignment = { vertical: 'center', horizontal: 'right' };
        }
      }
    }
  });

  // Section 2 Header (Slow Moving)
  const slowHeaderRef = `A${slowHeaderRowIdx + 1}`;
  if (ws1[slowHeaderRef]) {
    ws1[slowHeaderRef].s = { fill: { fgColor: { rgb: 'E11D48' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 11 }, alignment: { vertical: 'center' } };
  }
  for (let c = 0; c < 7; c++) {
    const colName = String.fromCharCode(65 + c);
    const cellRef = `${colName}${slowColsRowIdx + 1}`;
    if (ws1[cellRef]) {
      ws1[cellRef].s = { fill: { fgColor: { rgb: 'BE123C' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 9 }, alignment: { vertical: 'center' } };
    }
  }

  // Section 2 Data rows
  leastList.forEach((_, idx) => {
    const r = slowColsRowIdx + 2 + idx;
    const bg = idx % 2 === 0 ? 'FFF1F2' : 'FFFFFF';
    for (let c = 0; c < 7; c++) {
      const cellRef = `${String.fromCharCode(65 + c)}${r}`;
      if (ws1[cellRef]) {
        ws1[cellRef].s = { fill: { fgColor: { rgb: bg } }, font: { sz: 9 }, alignment: { vertical: 'center' } };
        if (c === 4) {
          ws1[cellRef].z = '"Rp" #,##0';
          ws1[cellRef].s.alignment = { vertical: 'center', horizontal: 'right' };
        }
      }
    }
  });

  XLSX.utils.book_append_sheet(wb, ws1, 'Ringkasan & Top Items');


  // ==========================================
  // SHEET 2: Rincian Barang & Kategori
  // ==========================================
  const sheet2Data: any[][] = [
    ['LIST SELURUH BARANG LAKU & RINCIAN KATEGORI'],
    [`Rincian lengkap kinerja per barang, penjualan per merk, dan penjualan per tipe produk — Cetak: ${printedDate}`],
    [''], // Spacer

    // Section 1: All Sold Products Table
    ['📋 LIST SELURUH BARANG LAKU'],
    ['No', 'Nama Produk', 'Merk', 'Tipe Produk', 'Harga Satuan', 'Total Terjual', 'Total Sales (Omset)']
  ];

  const allSold = stats.all_sold_products || [];
  let totalAllSales = 0;
  let totalAllQty = 0;

  allSold.forEach((p, idx) => {
    totalAllSales += (p.total_sales || 0);
    totalAllQty += (p.total_qty || 0);
    sheet2Data.push([
      idx + 1,
      p.product_name,
      p.artist || 'Umum',
      p.product_type || 'Umum',
      p.price || 0,
      `${formatQuantity(p.total_qty)} ${unitLabel(p.unit)}`,
      p.total_sales || 0
    ]);
  });

  // Table 1 Summary Row
  const summaryRowIdx = sheet2Data.length;
  sheet2Data.push([
    'TOTAL',
    `${allSold.length} Jenis Produk`,
    '-',
    '-',
    '-',
    `${formatQuantity(totalAllQty)} Total Item`,
    totalAllSales
  ]);

  // Section 2: Penjualan Per Merk
  sheet2Data.push(['']);
  const merkHeaderRowIdx = sheet2Data.length;
  sheet2Data.push(['🏷️ RINGKASAN PENJUALAN PER MERK']);
  const merkColsRowIdx = sheet2Data.length;
  sheet2Data.push(['No', 'Nama Merk', 'Total Kuantitas Terjual', 'Total Sales (Omset)']);

  const artistSales = stats.sales_by_artist || [];
  artistSales.forEach((a, idx) => {
    sheet2Data.push([
      idx + 1,
      a.name,
      `${formatQuantity(a.total_qty)} Kuantitas`,
      a.total_sales || 0
    ]);
  });

  // Section 3: Penjualan Per Tipe Produk
  sheet2Data.push(['']);
  const typeHeaderRowIdx = sheet2Data.length;
  sheet2Data.push(['🔲 RINGKASAN PENJUALAN PER TIPE PRODUK']);
  const typeColsRowIdx = sheet2Data.length;
  sheet2Data.push(['No', 'Tipe Produk', 'Total Kuantitas Terjual', 'Total Sales (Omset)']);

  const typeSales = stats.sales_by_type || [];
  typeSales.forEach((t, idx) => {
    sheet2Data.push([
      idx + 1,
      t.name,
      `${formatQuantity(t.total_qty)} Kuantitas`,
      t.total_sales || 0
    ]);
  });

  const ws2 = XLSX.utils.aoa_to_sheet(sheet2Data);

  // Set Merges for Sheet 2
  ws2['!merges'] = [
    { s: { r: 0, c: 0 }, e: { r: 0, c: 6 } }, // Title
    { s: { r: 1, c: 0 }, e: { r: 1, c: 6 } }, // Subtitle
    { s: { r: 3, c: 0 }, e: { r: 3, c: 6 } }, // All Sold Header
    { s: { r: merkHeaderRowIdx, c: 0 }, e: { r: merkHeaderRowIdx, c: 3 } }, // Merk Header
    { s: { r: typeHeaderRowIdx, c: 0 }, e: { r: typeHeaderRowIdx, c: 3 } }  // Type Header
  ];

  ws2['!cols'] = [
    { wch: 8 },  // No
    { wch: 36 }, // Nama Produk / Merk / Tipe
    { wch: 22 }, // Merk / Qty
    { wch: 22 }, // Tipe Produk / Sales
    { wch: 18 }, // Harga Satuan
    { wch: 22 }, // Total Terjual
    { wch: 24 }  // Total Sales
  ];

  // Apply Styling for Sheet 2
  if (ws2['A1']) ws2['A1'].s = { fill: { fgColor: { rgb: '1E293B' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 14 }, alignment: { vertical: 'center' } };
  if (ws2['A2']) ws2['A2'].s = { fill: { fgColor: { rgb: '1E293B' } }, font: { color: { rgb: '94A3B8' }, sz: 9, italic: true }, alignment: { vertical: 'center' } };

  // All Sold Header
  if (ws2['A4']) ws2['A4'].s = { fill: { fgColor: { rgb: '334155' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 11 }, alignment: { vertical: 'center' } };
  for (let c = 0; c < 7; c++) {
    const colName = String.fromCharCode(65 + c);
    if (ws2[`${colName}5`]) {
      ws2[`${colName}5`].s = { fill: { fgColor: { rgb: '0F172A' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 9 }, alignment: { vertical: 'center' } };
    }
  }

  // All Sold Data Rows
  allSold.forEach((_, idx) => {
    const r = 6 + idx;
    const bg = idx % 2 === 0 ? 'F8FAFC' : 'FFFFFF';
    for (let c = 0; c < 7; c++) {
      const cellRef = `${String.fromCharCode(65 + c)}${r}`;
      if (ws2[cellRef]) {
        ws2[cellRef].s = { fill: { fgColor: { rgb: bg } }, font: { sz: 9 }, alignment: { vertical: 'center' } };
        if (c === 4 || c === 6) {
          ws2[cellRef].z = '"Rp" #,##0';
          ws2[cellRef].s.alignment = { vertical: 'center', horizontal: 'right' };
        }
      }
    }
  });

  // Table 1 Summary Row Style
  const sumRowRefIdx = summaryRowIdx + 1;
  for (let c = 0; c < 7; c++) {
    const cellRef = `${String.fromCharCode(65 + c)}${sumRowRefIdx}`;
    if (ws2[cellRef]) {
      ws2[cellRef].s = { fill: { fgColor: { rgb: 'E2E8F0' } }, font: { bold: true, color: { rgb: '0F172A' }, sz: 9 }, alignment: { vertical: 'center' } };
      if (c === 6) {
        ws2[cellRef].z = '"Rp" #,##0';
        ws2[cellRef].s.alignment = { vertical: 'center', horizontal: 'right' };
      }
    }
  }

  // Merk Header Style
  const merkHeaderRef = `A${merkHeaderRowIdx + 1}`;
  if (ws2[merkHeaderRef]) {
    ws2[merkHeaderRef].s = { fill: { fgColor: { rgb: 'D97706' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 10 }, alignment: { vertical: 'center' } };
  }
  for (let c = 0; c < 4; c++) {
    const colName = String.fromCharCode(65 + c);
    const cellRef = `${colName}${merkColsRowIdx + 1}`;
    if (ws2[cellRef]) {
      ws2[cellRef].s = { fill: { fgColor: { rgb: 'B45309' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 9 }, alignment: { vertical: 'center' } };
    }
  }

  artistSales.forEach((_, idx) => {
    const r = merkColsRowIdx + 2 + idx;
    const bg = idx % 2 === 0 ? 'FFFBEB' : 'FFFFFF';
    for (let c = 0; c < 4; c++) {
      const cellRef = `${String.fromCharCode(65 + c)}${r}`;
      if (ws2[cellRef]) {
        ws2[cellRef].s = { fill: { fgColor: { rgb: bg } }, font: { sz: 9 }, alignment: { vertical: 'center' } };
        if (c === 3) {
          ws2[cellRef].z = '"Rp" #,##0';
          ws2[cellRef].s.alignment = { vertical: 'center', horizontal: 'right' };
        }
      }
    }
  });

  // Type Header Style
  const typeHeaderRef = `A${typeHeaderRowIdx + 1}`;
  if (ws2[typeHeaderRef]) {
    ws2[typeHeaderRef].s = { fill: { fgColor: { rgb: 'DB2777' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 10 }, alignment: { vertical: 'center' } };
  }
  for (let c = 0; c < 4; c++) {
    const colName = String.fromCharCode(65 + c);
    const cellRef = `${colName}${typeColsRowIdx + 1}`;
    if (ws2[cellRef]) {
      ws2[cellRef].s = { fill: { fgColor: { rgb: 'BE185D' } }, font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 9 }, alignment: { vertical: 'center' } };
    }
  }

  typeSales.forEach((_, idx) => {
    const r = typeColsRowIdx + 2 + idx;
    const bg = idx % 2 === 0 ? 'FDF2F8' : 'FFFFFF';
    for (let c = 0; c < 4; c++) {
      const cellRef = `${String.fromCharCode(65 + c)}${r}`;
      if (ws2[cellRef]) {
        ws2[cellRef].s = { fill: { fgColor: { rgb: bg } }, font: { sz: 9 }, alignment: { vertical: 'center' } };
        if (c === 3) {
          ws2[cellRef].z = '"Rp" #,##0';
          ws2[cellRef].s.alignment = { vertical: 'center', horizontal: 'right' };
        }
      }
    }
  });

  XLSX.utils.book_append_sheet(wb, ws2, 'Detail Barang & Kategori');

  // Generate File Output
  const fileName = `Laporan_Penjualan_${new Date().toISOString().slice(0, 10)}.xlsx`;
  const wbout = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
  const blob = new Blob([wbout], { type: 'application/octet-stream' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
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
  doc.text(`${stats.total_items_sold}`, 142, startY + 14);

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
    `${p.total_qty} ${unitLabel(p.unit)}`,
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
    `${p.total_qty} ${unitLabel(p.unit)}`,
    `Sisa ${p.stock ?? 0} ${unitLabel(p.unit)}`
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
    `${p.total_qty} ${unitLabel(p.unit)}`,
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

export const exportToCSV = (stats: DashboardStats, dateTitle: string = 'Keseluruhan'): void => {
  const allSold = stats.all_sold_products || [];
  
  const headers = ['No', 'Nama Produk', 'Merk/Artist', 'Tipe Produk', 'Harga Satuan', 'Kuantitas Terjual', 'Total Omset'];
  const csvRows = [headers.join(',')];
  
  allSold.forEach((p, idx) => {
    const row = [
      idx + 1,
      `"${(p.product_name || '').replace(/"/g, '""')}"`,
      `"${(p.artist || 'Umum').replace(/"/g, '""')}"`,
      `"${(p.product_type || 'Umum').replace(/"/g, '""')}"`,
      p.price || 0,
      p.total_qty || 0,
      p.total_sales || 0
    ];
    csvRows.push(row.join(','));
  });

  const csvString = csvRows.join('\n');
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
