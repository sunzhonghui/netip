/**
 * Domain resolution utilities for probe subdomains and API hosts.
 */

export interface DomainInfo {
  baseDomain: string
  apiHost: string
}

const DEFAULT_BASE_DOMAIN = 'ipw.3x.cx'
const DEFAULT_API_HOST = 'ip.ipw.3x.cx'

/**
 * Resolves the base domain (e.g. 'ipw.3x.cx') and API host (e.g. 'ip.ipw.3x.cx')
 * based on current browser hostname or given hostname.
 */
export function getDomainInfo(hostname?: string): DomainInfo {
  const host = (hostname ?? (typeof window !== 'undefined' ? window.location.hostname : '')).toLowerCase().trim()

  const isLocalOrIP =
    !host ||
    host === 'localhost' ||
    host === '127.0.0.1' ||
    host === '::1' ||
    !host.includes('.') ||
    /^(\d{1,3}\.){3}\d{1,3}$/.test(host) ||
    host.includes(':')

  // In local dev, preview, or direct IP access, default to ipw.3x.cx
  if (isLocalOrIP) {
    return {
      baseDomain: DEFAULT_BASE_DOMAIN,
      apiHost: DEFAULT_API_HOST,
    }
  }

  // If host is ipw.3x.cx or any subdomain of ipw.3x.cx (e.g. ip.ipw.3x.cx, 4.ipw.3x.cx, test.ipw.3x.cx)
  if (host === 'ipw.3x.cx' || host.endsWith('.ipw.3x.cx')) {
    return {
      baseDomain: 'ipw.3x.cx',
      apiHost: host.startsWith('ip.') ? host : 'ip.ipw.3x.cx',
    }
  }

  // Generic custom domain support:
  // If host starts with a recognized probe or web prefix (ip., 4., 6., test.), strip it
  const parts = host.split('.')
  if (parts.length > 2 && ['ip', '4', '6', 'test'].includes(parts[0])) {
    const baseDomain = parts.slice(1).join('.')
    return {
      baseDomain,
      apiHost: parts[0] === 'ip' ? host : `ip.${baseDomain}`,
    }
  }

  // Apex or standard domain (e.g. example.com)
  return {
    baseDomain: host,
    apiHost: `ip.${host}`,
  }
}
