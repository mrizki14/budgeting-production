const months = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
]

export function formatMoney(amount) {
  const value = Number(amount ?? 0)
  return `Rp ${new Intl.NumberFormat('id-ID', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value)}`
}

export function formatDate(value) {
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: '2-digit',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(value))
}

export function monthName(month) {
  return months[Number(month) - 1] ?? ''
}

export function dateInputValue(value) {
  return value ? String(value).slice(0, 10) : ''
}
