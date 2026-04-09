function stripNonDigits(value: string) {
  return value.replace(/\D+/g, '')
}

function formatCurrencyDigits(digits: string) {
  if (!digits) return ''

  const padded = digits.padStart(3, '0')
  const cents = padded.slice(-2)
  const reais = padded.slice(0, -2).replace(/^0+(?=\d)/, '') || '0'

  return `${Number(reais).toLocaleString('pt-BR')},${cents}`
}

export function normalizeBrazilianPhone(value: string) {
  return stripNonDigits(value).slice(0, 11)
}

export function formatBrazilianPhoneInput(value: string) {
  const digits = normalizeBrazilianPhone(value)
  if (!digits) return ''

  if (digits.length <= 2) return `(${digits}`
  if (digits.length <= 6) return `(${digits.slice(0, 2)}) ${digits.slice(2)}`
  if (digits.length <= 10) return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`
  return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`
}

export function formatBrazilianPhoneDisplay(value: string | null | undefined) {
  if (!value) return '-'
  return formatBrazilianPhoneInput(value)
}

export function normalizeBrazilianCurrencyInput(value: string) {
  const digits = stripNonDigits(value)
  if (!digits) return undefined
  return Number(digits)
}

export function formatBrazilianCurrencyInput(value: string) {
  const digits = stripNonDigits(value)
  return formatCurrencyDigits(digits)
}

export function formatBrazilianCurrencyFromCents(cents: number | null | undefined) {
  if (cents == null) return ''
  return formatCurrencyDigits(String(Math.max(0, Math.trunc(cents))))
}

