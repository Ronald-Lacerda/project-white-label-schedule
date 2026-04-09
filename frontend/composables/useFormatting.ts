function pad(value: number) {
  return String(value).padStart(2, '0')
}

export function onlyDigits(value: string) {
  return value.replace(/\D/g, '')
}

export function normalizePhoneBr(value: string) {
  return onlyDigits(value).slice(0, 11)
}

export function formatPhoneBr(value: string) {
  const digits = normalizePhoneBr(value)
  if (!digits) return ''

  if (digits.length <= 2) {
    return `(${digits}`
  }

  if (digits.length <= 6) {
    return `(${digits.slice(0, 2)}) ${digits.slice(2)}`
  }

  if (digits.length <= 10) {
    return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`
  }

  return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`
}

export function formatCentsToMoneyInput(cents: number | null | undefined) {
  if (cents == null) return ''
  return (cents / 100).toLocaleString('pt-BR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

export function formatMoneyInputBr(value: string) {
  const digits = onlyDigits(value)
  if (!digits) return ''
  return (Number(digits) / 100).toLocaleString('pt-BR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

export function parseMoneyInputToCents(value: string) {
  const digits = onlyDigits(value)
  if (!digits) return null
  return Number(digits)
}

export function formatShortTimeBr(value: string) {
  if (!value) return ''

  if (/^\d{1,2}:\d{2}/.test(value)) {
    const [hoursRaw, minutesRaw = '00'] = value.split(':')
    const hours = Number(hoursRaw)
    const minutes = Number(minutesRaw)
    return minutes === 0 ? `${hours}h` : `${hours}h${pad(minutes)}`
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value

  const hours = date.getHours()
  const minutes = date.getMinutes()
  return minutes === 0 ? `${hours}h` : `${hours}h${pad(minutes)}`
}

export function formatTimeRangeBr(start: string, end: string) {
  return `${formatShortTimeBr(start)} as ${formatShortTimeBr(end)}`
}
