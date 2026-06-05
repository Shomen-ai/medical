// Файл: composables/useExcel.ts
// Назначение: генерация .xlsx (SheetJS) из листов с форматированием (ширины колонок,
// числовые форматы валюты/десятичные, объединённый заголовок). xlsx подгружается
// динамически — только на клиенте, вне SSR-бандла.

// Один лист отчёта.
export interface ReportSheet {
  name: string                       // имя вкладки (Excel обрежет до 31 символа)
  rows: (string | number)[][]        // строки (массив массивов)
  cols?: number[]                    // ширины колонок в символах
  money?: number[]                   // индексы колонок с денежным форматом (# ##0)
  decimal?: number[]                 // индексы колонок с форматом 0.00
  mergeTitleRow?: boolean            // объединить первую строку (заголовок) по всем колонкам
}

// Собирает книгу .xlsx из листов и инициирует скачивание файла (вызывать на клиенте).
export async function downloadXlsx(filename: string, sheets: ReportSheet[]): Promise<void> {
  const XLSX = await import('xlsx')
  const wb = XLSX.utils.book_new()

  const applyFmt = (ws: any, r: number, c: number, z: string) => {
    const addr = XLSX.utils.encode_cell({ r, c })
    const cell = ws[addr]
    if (cell && typeof cell.v === 'number') cell.z = z
  }

  for (const s of sheets) {
    const ws = XLSX.utils.aoa_to_sheet(s.rows)
    if (s.cols) ws['!cols'] = s.cols.map(w => ({ wch: w }))

    const ref = ws['!ref']
    if (ref) {
      const range = XLSX.utils.decode_range(ref)
      for (let r = range.s.r; r <= range.e.r; r++) {
        for (const c of s.money ?? []) applyFmt(ws, r, c, '#,##0')
        for (const c of s.decimal ?? []) applyFmt(ws, r, c, '0.00')
      }
      if (s.mergeTitleRow) {
        ws['!merges'] = [{ s: { r: 0, c: 0 }, e: { r: 0, c: range.e.c } }]
      }
    }

    XLSX.utils.book_append_sheet(wb, ws, s.name.slice(0, 31))
  }
  XLSX.writeFile(wb, filename)
}
