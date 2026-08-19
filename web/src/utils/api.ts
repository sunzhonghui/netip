import type {
  APIResponse,
  IPDetails,
  DNSQueryResult,
  MultiNodePingResponse,
  MultiNodeTCPingResponse,
  IPv6CheckResponse,
  SSLCheckResult,
  WHOISResult,
  ASNResult,
  HTTPCheckResult,
  SpeedTestResult,
} from '@/types'

const API_BASE = '/api/v1'

export async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const url = path.startsWith('http') ? path : `${API_BASE}${path}`
  const headers = new Headers(options.headers || {})
  if (!headers.has('Content-Type') && options.method && options.method !== 'GET') {
    headers.set('Content-Type', 'application/json')
  }

  const resp = await fetch(url, {
    ...options,
    headers,
  })

  const json: APIResponse<T> = await resp.json()
  if (!json.success || !json.data) {
    const msg = json.error?.message || `请求失败 (${resp.status})`
    throw new Error(msg)
  }

  return json.data
}

// API methods
export const api = {
  getMe(): Promise<IPDetails> {
    return request<IPDetails>('/me')
  },

  getIP(ip: string): Promise<IPDetails> {
    return request<IPDetails>(`/ip/${encodeURIComponent(ip)}`)
  },

  getASN(query: string): Promise<ASNResult> {
    return request<ASNResult>(`/asn/${encodeURIComponent(query)}`)
  },

  queryDNS(name: string, type: string): Promise<DNSQueryResult> {
    return request<DNSQueryResult>('/dns', {
      method: 'POST',
      body: JSON.stringify({ name, type }),
    })
  },

  ping(target: string, count = 4, ipVersion = 'auto'): Promise<MultiNodePingResponse> {
    return request<MultiNodePingResponse>('/ping', {
      method: 'POST',
      body: JSON.stringify({ target, count, ip_version: ipVersion }),
    })
  },

  tcping(target: string, port = 80, count = 4): Promise<MultiNodeTCPingResponse> {
    return request<MultiNodeTCPingResponse>('/tcping', {
      method: 'POST',
      body: JSON.stringify({ target, port, count }),
    })
  },

  checkIPv6(target: string): Promise<IPv6CheckResponse> {
    return request<IPv6CheckResponse>('/ipv6-check', {
      method: 'POST',
      body: JSON.stringify({ target }),
    })
  },

  checkSSL(hostname: string, port = 443): Promise<SSLCheckResult> {
    return request<SSLCheckResult>('/ssl', {
      method: 'POST',
      body: JSON.stringify({ hostname, port }),
    })
  },

  checkHTTP(target: string): Promise<HTTPCheckResult> {
    return request<HTTPCheckResult>('/http', {
      method: 'POST',
      body: JSON.stringify({ target }),
    })
  },

  speedTest(target: string): Promise<SpeedTestResult> {
    return request<SpeedTestResult>('/speed', {
      method: 'POST',
      body: JSON.stringify({ target }),
    })
  },

  queryWHOIS(target: string): Promise<WHOISResult> {
    return request<WHOISResult>('/whois', {
      method: 'POST',
      body: JSON.stringify({ target }),
    })
  },
}

// Fetch with timeout for homepage IPv4 / IPv6 detection
export async function fetchWithTimeout<T>(url: string, timeoutMs = 4000): Promise<T> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs)

  try {
    const resp = await fetch(url, {
      signal: controller.signal,
      headers: {
        Accept: 'application/json',
      },
    })
    clearTimeout(timeoutId)
    if (!resp.ok) {
      throw new Error(`HTTP ${resp.status}`)
    }
    return await resp.json()
  } catch (err) {
    clearTimeout(timeoutId)
    throw err
  }
}
