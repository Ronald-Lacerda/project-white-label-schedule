function pad(value: number) {
  return String(value).padStart(2, '0')
}

export function formatLocalDateInputValue(date = new Date()) {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

export function formatLocalDateTimeInputValue(date = new Date()) {
  return `${formatLocalDateInputValue(date)}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export function localDateTimeToIso(value: string) {
  return new Date(value).toISOString()
}
