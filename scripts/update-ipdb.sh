#!/usr/bin/env bash
# ==============================================================================
# NetIP - Offline IP Database Auto-Update Script
# Downloads / updates ip2region.xdb, GeoLite2-ASN.mmdb, GeoLite2-City.mmdb
# into the persistent ./data/ipdb directory.
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
IPDB_DIR="${ROOT_DIR}/data/ipdb"

mkdir -p "${IPDB_DIR}"
echo "[+] Target IP Database Directory: ${IPDB_DIR}"

# 1. Update ip2region.xdb
echo "[+] Downloading ip2region.xdb..."
IP2REGION_URLS=(
    "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb"
    "https://fastly.jsdelivr.net/gh/lionsoul2014/ip2region@master/data/ip2region.xdb"
    "https://raw.gitmirror.com/lionsoul2014/ip2region/master/data/ip2region.xdb"
)

downloaded_ip2region=0
for url in "${IP2REGION_URLS[@]}"; do
    echo "    Trying mirror: ${url}"
    if curl -sSL --connect-timeout 10 --max-time 60 "${url}" -o "${IPDB_DIR}/ip2region.xdb.tmp"; then
        if [ -s "${IPDB_DIR}/ip2region.xdb.tmp" ] && [ $(stat -c%s "${IPDB_DIR}/ip2region.xdb.tmp" 2>/dev/null || stat -f%z "${IPDB_DIR}/ip2region.xdb.tmp") -gt 1000000 ]; then
            mv "${IPDB_DIR}/ip2region.xdb.tmp" "${IPDB_DIR}/ip2region.xdb"
            echo "    ✓ ip2region.xdb updated successfully."
            downloaded_ip2region=1
            break
        fi
    fi
    rm -f "${IPDB_DIR}/ip2region.xdb.tmp"
done

# 2. Update GeoLite2-ASN.mmdb
echo "[+] Downloading GeoLite2-ASN.mmdb..."
ASN_URLS=(
    "https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-ASN.mmdb"
    "https://git.io/GeoLite2-ASN.mmdb"
)

for url in "${ASN_URLS[@]}"; do
    echo "    Trying mirror: ${url}"
    if curl -sSL --connect-timeout 10 --max-time 120 "${url}" -o "${IPDB_DIR}/GeoLite2-ASN.mmdb.tmp"; then
        if [ -s "${IPDB_DIR}/GeoLite2-ASN.mmdb.tmp" ] && [ $(stat -c%s "${IPDB_DIR}/GeoLite2-ASN.mmdb.tmp" 2>/dev/null || stat -f%z "${IPDB_DIR}/GeoLite2-ASN.mmdb.tmp") -gt 2000000 ]; then
            mv "${IPDB_DIR}/GeoLite2-ASN.mmdb.tmp" "${IPDB_DIR}/GeoLite2-ASN.mmdb"
            echo "    ✓ GeoLite2-ASN.mmdb updated successfully."
            break
        fi
    fi
    rm -f "${IPDB_DIR}/GeoLite2-ASN.mmdb.tmp"
done

# 3. Update GeoLite2-City.mmdb
echo "[+] Downloading GeoLite2-City.mmdb..."
CITY_URLS=(
    "https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-City.mmdb"
    "https://git.io/GeoLite2-City.mmdb"
)

for url in "${CITY_URLS[@]}"; do
    echo "    Trying mirror: ${url}"
    if curl -sSL --connect-timeout 10 --max-time 180 "${url}" -o "${IPDB_DIR}/GeoLite2-City.mmdb.tmp"; then
        if [ -s "${IPDB_DIR}/GeoLite2-City.mmdb.tmp" ] && [ $(stat -c%s "${IPDB_DIR}/GeoLite2-City.mmdb.tmp" 2>/dev/null || stat -f%z "${IPDB_DIR}/GeoLite2-City.mmdb.tmp") -gt 10000000 ]; then
            mv "${IPDB_DIR}/GeoLite2-City.mmdb.tmp" "${IPDB_DIR}/GeoLite2-City.mmdb"
            echo "    ✓ GeoLite2-City.mmdb updated successfully."
            break
        fi
    fi
    rm -f "${IPDB_DIR}/GeoLite2-City.mmdb.tmp"
done

echo "[+] Done! IP databases in ${IPDB_DIR}:"
ls -lh "${IPDB_DIR}"
