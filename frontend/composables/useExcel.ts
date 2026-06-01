// Файл: composables/useExcel.ts
// Назначение: генерация .xlsx (SheetJS) из простых листов и скачивание в браузере.
// xlsx подгружается динамически — только на клиенте, вне SSR-бандла.

// Один лист отчёта: имя вкладки + строки (массив массивов значений).
export interface ReportSheet {
  name: string
  rows: (string | number)[][]
}

// Собирает книгу .xlsx из листов и инициирует скачивание файла (вызывать на клиенте).
export async function downloadXlsx(filename: string, sheets: ReportSheet[]): Promise<void> {
  const XLSX = await import('xlsx')
  const wb = XLSX.utils.book_new()
  for (const s of sheets) {
    const ws = XLSX.utils.aoa_to_sheet(s.rows)
    // Имя вкладки в Excel ограничено 31 символом.
    XLSX.utils.book_append_sheet(wb, ws, s.name.slice(0, 31))
  }
  XLSX.writeFile(wb, filename)
}
