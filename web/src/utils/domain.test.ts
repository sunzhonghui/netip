import { describe, it, expect } from 'vitest'
import { getDomainInfo } from './domain'

describe('getDomainInfo', () => {
  it('handles apex domain ipw.3x.cx', () => {
    const info = getDomainInfo('ipw.3x.cx')
    expect(info.baseDomain).toBe('ipw.3x.cx')
    expect(info.apiHost).toBe('ip.ipw.3x.cx')
  })

  it('handles ip.ipw.3x.cx', () => {
    const info = getDomainInfo('ip.ipw.3x.cx')
    expect(info.baseDomain).toBe('ipw.3x.cx')
    expect(info.apiHost).toBe('ip.ipw.3x.cx')
  })

  it('handles 4.ipw.3x.cx', () => {
    const info = getDomainInfo('4.ipw.3x.cx')
    expect(info.baseDomain).toBe('ipw.3x.cx')
    expect(info.apiHost).toBe('ip.ipw.3x.cx')
  })

  it('handles 6.ipw.3x.cx', () => {
    const info = getDomainInfo('6.ipw.3x.cx')
    expect(info.baseDomain).toBe('ipw.3x.cx')
    expect(info.apiHost).toBe('ip.ipw.3x.cx')
  })

  it('handles test.ipw.3x.cx', () => {
    const info = getDomainInfo('test.ipw.3x.cx')
    expect(info.baseDomain).toBe('ipw.3x.cx')
    expect(info.apiHost).toBe('ip.ipw.3x.cx')
  })

  it('falls back to ipw.3x.cx for localhost and loopback', () => {
    expect(getDomainInfo('localhost')).toEqual({
      baseDomain: 'ipw.3x.cx',
      apiHost: 'ip.ipw.3x.cx',
    })
    expect(getDomainInfo('127.0.0.1')).toEqual({
      baseDomain: 'ipw.3x.cx',
      apiHost: 'ip.ipw.3x.cx',
    })
    expect(getDomainInfo('::1')).toEqual({
      baseDomain: 'ipw.3x.cx',
      apiHost: 'ip.ipw.3x.cx',
    })
  })

  it('falls back to ipw.3x.cx for IP addresses', () => {
    expect(getDomainInfo('192.168.1.100')).toEqual({
      baseDomain: 'ipw.3x.cx',
      apiHost: 'ip.ipw.3x.cx',
    })
  })

  it('handles custom domains', () => {
    expect(getDomainInfo('example.com')).toEqual({
      baseDomain: 'example.com',
      apiHost: 'ip.example.com',
    })
    expect(getDomainInfo('ip.example.com')).toEqual({
      baseDomain: 'example.com',
      apiHost: 'ip.example.com',
    })
    expect(getDomainInfo('4.example.com')).toEqual({
      baseDomain: 'example.com',
      apiHost: 'ip.example.com',
    })
  })
})
