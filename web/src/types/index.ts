export interface APIResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
  }
  request_id: string
}

export interface IPDetails {
  ip: string
  version: number
  country?: string
  country_code?: string
  province?: string
  city?: string
  isp?: string
  asn?: number
  as_name?: string
  network?: string
  sources?: Record<string, string>
}

export interface DualStackDetectionState {
  ipv4?: {
    ip: string
    details?: IPDetails
  }
  ipv6?: {
    ip: string
    details?: IPDetails
  }
  preferred?: 'ipv4' | 'ipv6' | 'unknown'
  loading: boolean
}

export interface DNSAnswer {
  type: string
  value: string
  ttl: number
  priority?: number
}

export interface ResolverResult {
  node_id?: string
  node_name?: string
  isp?: string
  resolver: string
  latency_ms: number
  answers: DNSAnswer[]
  error?: string
}

export interface DNSQueryResult {
  name: string
  type: string
  results: ResolverResult[]
}

export interface PingResult {
  target: string
  resolved_ip: string
  node: string
  sent: number
  received: number
  loss_percent: number
  min_ms: number
  avg_ms: number
  max_ms: number
  samples: number[]
  error?: string
}

export interface MultiNodePingResponse {
  target: string
  resolved_ip: string
  nodes: PingResult[]
}

export interface TCPingSample {
  success: boolean
  latency_ms?: number
  error?: string
}

export interface TCPingResult {
  target: string
  port: number
  resolved_ip: string
  node: string
  success: number
  failed: number
  avg_ms: number
  min_ms: number
  max_ms: number
  samples: TCPingSample[]
}

export interface MultiNodeTCPingResponse {
  target: string
  port: number
  resolved_ip: string
  nodes: TCPingResult[]
}

export interface EndpointStatus {
  supported: boolean
  status_code?: number
  latency_ms?: number
  error?: string
}

export interface ProtocolCheckResult {
  dns: boolean
  addresses: string[]
  http: EndpointStatus
  https: EndpointStatus
}

export interface IPv6CheckResponse {
  domain: string
  ipv4: ProtocolCheckResult
  ipv6: ProtocolCheckResult
  supported: boolean
  conclusion: string
}

export interface SSLCertificateInfo {
  subject: string
  issuer: string
  serial_number: string
  dns_names: string[]
  not_before: string
  not_after: string
  days_remaining: number
}

export interface SSLCheckResult {
  hostname: string
  port: number
  resolved_ip: string
  valid: boolean
  days_remaining: number
  issuer: string
  subject: string
  dns_names: string[]
  not_before: string
  not_after: string
  tls_version: string
  cipher_suite: string
  certificates: SSLCertificateInfo[]
  error?: string
}

export interface WHOISResult {
  query: string
  type: string
  domain?: string
  registrar?: string
  created?: string
  updated?: string
  expires?: string
  status?: string[]
  name_servers?: string[]
  dnssec?: string
  network?: string
  cidr?: string
  organization?: string
  country?: string
  asn?: string
  raw_text?: string
  source: string
}

export interface ASNResult {
  asn: number
  as_name: string
  country: string
  registry?: string
  network?: string
  source: string
}

export interface HTTPCheckResult {
  url: string
  dns_ms: number
  tcp_ms: number
  tls_ms: number
  ttfb_ms: number
  total_ms: number
  protocol: string
  status_code: number
  status_text: string
  resolved_ip: string
  server?: string
  content_type?: string
  content_length?: number
  headers?: Record<string, string>
}

export interface SpeedTestResult {
  target: string
  dns_ms: number
  connect_ms: number
  tls_ms: number
  ttfb_ms: number
  download_bytes: number
  download_ms: number
  speed_mbps: number
  resolved_ip: string
}
