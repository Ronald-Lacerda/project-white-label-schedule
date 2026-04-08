type NullableColor = string | null | undefined

export interface BrandThemeInput {
  primaryColor?: NullableColor
  secondaryColor?: NullableColor
}

const DEFAULT_PRIMARY = '#1f3b73'
const DEFAULT_SECONDARY = '#c99249'

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max)
}

function normalizeHex(color: NullableColor, fallback: string) {
  const raw = (color ?? '').trim()
  const value = raw.startsWith('#') ? raw.slice(1) : raw

  if (/^[0-9a-fA-F]{3}$/.test(value)) {
    return `#${value.split('').map(char => `${char}${char}`).join('').toLowerCase()}`
  }

  if (/^[0-9a-fA-F]{6}$/.test(value)) {
    return `#${value.toLowerCase()}`
  }

  return fallback
}

function hexToRgb(hex: string) {
  const value = normalizeHex(hex, DEFAULT_PRIMARY).slice(1)
  return {
    r: Number.parseInt(value.slice(0, 2), 16),
    g: Number.parseInt(value.slice(2, 4), 16),
    b: Number.parseInt(value.slice(4, 6), 16),
  }
}

function rgbToHex({ r, g, b }: { r: number; g: number; b: number }) {
  const toHex = (channel: number) => clamp(Math.round(channel), 0, 255).toString(16).padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

function mixColors(colorA: string, colorB: string, weightA: number) {
  const a = hexToRgb(colorA)
  const b = hexToRgb(colorB)
  const mix = (first: number, second: number) => first * weightA + second * (1 - weightA)
  return rgbToHex({
    r: mix(a.r, b.r),
    g: mix(a.g, b.g),
    b: mix(a.b, b.b),
  })
}

function toChannelString(color: string) {
  const rgb = hexToRgb(color)
  return `${rgb.r} ${rgb.g} ${rgb.b}`
}

function luminanceChannel(channel: number) {
  const normalized = channel / 255
  return normalized <= 0.03928 ? normalized / 12.92 : ((normalized + 0.055) / 1.055) ** 2.4
}

function relativeLuminance(color: string) {
  const { r, g, b } = hexToRgb(color)
  return 0.2126 * luminanceChannel(r) + 0.7152 * luminanceChannel(g) + 0.0722 * luminanceChannel(b)
}

export function contrastRatio(first: string, second: string) {
  const luminanceA = relativeLuminance(first)
  const luminanceB = relativeLuminance(second)
  const lighter = Math.max(luminanceA, luminanceB)
  const darker = Math.min(luminanceA, luminanceB)
  return (lighter + 0.05) / (darker + 0.05)
}

export function readableTextColor(background: string) {
  const whiteContrast = contrastRatio(background, '#ffffff')
  const darkContrast = contrastRatio(background, '#172033')
  return whiteContrast >= darkContrast ? '#ffffff' : '#172033'
}

export function buildBrandTheme(input: BrandThemeInput = {}) {
  const primary = normalizeHex(input.primaryColor, DEFAULT_PRIMARY)
  const secondary = normalizeHex(input.secondaryColor, DEFAULT_SECONDARY)

  const brandSurface = mixColors(secondary, '#ffffff', 0.18)
  const brandSurfaceMuted = mixColors(primary, '#ffffff', 0.09)
  const brandAccent = mixColors(primary, '#ffffff', 0.18)

  return {
    primary,
    secondary,
    accent: brandAccent,
    surface: brandSurface,
    surfaceMuted: brandSurfaceMuted,
    text: '#172033',
    success: '#147d64',
    warning: '#b86b16',
    danger: '#bf3a36',
    info: '#2a65b6',
    onPrimary: readableTextColor(primary),
    onSecondary: readableTextColor(secondary),
    primaryRgb: toChannelString(primary),
    secondaryRgb: toChannelString(secondary),
  }
}

export function brandThemeStyle(input: BrandThemeInput = {}) {
  const theme = buildBrandTheme(input)

  return {
    '--color-brand-primary': theme.primary,
    '--color-brand-secondary': theme.secondary,
    '--color-brand-accent': theme.accent,
    '--color-brand-surface': theme.surface,
    '--color-brand-surface-muted': theme.surfaceMuted,
    '--color-brand-text': theme.text,
    '--color-brand-on-primary': theme.onPrimary,
    '--color-brand-on-secondary': theme.onSecondary,
    '--color-brand-primary-rgb': theme.primaryRgb,
    '--color-brand-secondary-rgb': theme.secondaryRgb,
  } as Record<string, string>
}

export function applyBrandTheme(input: BrandThemeInput = {}) {
  if (!import.meta.client) return

  const theme = brandThemeStyle(input)
  Object.entries(theme).forEach(([key, value]) => {
    document.documentElement.style.setProperty(key, value)
  })
}
