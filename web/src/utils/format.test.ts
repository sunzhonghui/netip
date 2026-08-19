import { describe, it, expect } from 'vitest'
import { getLatencyLevel, formatBytes, formatDate } from './format'

describe('format utilities', () => {
  it('getLatencyLevel should categorize response times correctly', () => {
    expect(getLatencyLevel(30).color).toBe('emerald')
    expect(getLatencyLevel(75).color).toBe('blue')
    expect(getLatencyLevel(150).color).toBe('amber')
    expect(getLatencyLevel(250).color).toBe('rose')
    expect(getLatencyLevel(-1).color).toBe('gray')
  })

  it('formatBytes should format byte quantities with units', () => {
    expect(formatBytes(0)).toBe('0 B')
    expect(formatBytes(1024)).toBe('1 KB')
    expect(formatBytes(1048576)).toBe('1 MB')
    expect(formatBytes(5242880)).toBe('5 MB')
  })

  it('formatDate should format ISO date strings', () => {
    expect(formatDate('')).toBe('-')
    const formatted = formatDate('2026-08-19T12:00:00Z')
    expect(formatted).not.toBe('-')
    expect(typeof formatted).toBe('string')
  })
})
