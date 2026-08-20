import{c as a}from"./index-CHOls667.js";/**
 * @license lucide-vue-next v0.475.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const y=a("ArrowRightIcon",[["path",{d:"M5 12h14",key:"1ays0h"}],["path",{d:"m12 5 7 7-7 7",key:"xquz4c"}]]);/**
 * @license lucide-vue-next v0.475.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const S=a("SearchIcon",[["circle",{cx:"11",cy:"11",r:"8",key:"4ej97u"}],["path",{d:"m21 21-4.3-4.3",key:"1qie3q"}]]),d="/api/v1";async function r(t,e={}){var i;const n=t.startsWith("http")?t:`${d}${t}`,s=new Headers(e.headers||{});!s.has("Content-Type")&&e.method&&e.method!=="GET"&&s.set("Content-Type","application/json");const o=await fetch(n,{...e,headers:s}),c=await o.json();if(!c.success||!c.data){const h=((i=c.error)==null?void 0:i.message)||`请求失败 (${o.status})`;throw new Error(h)}return c.data}const m={getMe(){return r("/me")},getIP(t){return r(`/ip/${encodeURIComponent(t)}`)},getASN(t){return r(`/asn/${encodeURIComponent(t)}`)},queryDNS(t,e){return r("/dns",{method:"POST",body:JSON.stringify({name:t,type:e})})},ping(t,e=4,n="auto"){return r("/ping",{method:"POST",body:JSON.stringify({target:t,count:e,ip_version:n})})},tcping(t,e=80,n=4){return r("/tcping",{method:"POST",body:JSON.stringify({target:t,port:e,count:n})})},checkIPv6(t){return r("/ipv6-check",{method:"POST",body:JSON.stringify({target:t})})},checkSSL(t,e=443){return r("/ssl",{method:"POST",body:JSON.stringify({hostname:t,port:e})})},checkHTTP(t){return r("/http",{method:"POST",body:JSON.stringify({target:t})})},speedTest(t){return r("/speed",{method:"POST",body:JSON.stringify({target:t})})},queryWHOIS(t){return r("/whois",{method:"POST",body:JSON.stringify({target:t})})}};async function p(t,e=4e3){const n=new AbortController,s=setTimeout(()=>n.abort(),e);try{const o=await fetch(t,{signal:n.signal,headers:{Accept:"application/json"}});if(clearTimeout(s),!o.ok)throw new Error(`HTTP ${o.status}`);return await o.json()}catch(o){throw clearTimeout(s),o}}export{y as A,S,m as a,p as f};
