import{c as d,d as i,a as l,n as u,q as m,g as n,u as s,t as b,f as h,l as p,o as a}from"./index-CHOls667.js";/**
 * @license lucide-vue-next v0.475.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const x=d("CheckIcon",[["path",{d:"M20 6 9 17l-5-5",key:"1gmf2c"}]]);/**
 * @license lucide-vue-next v0.475.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const y=d("CopyIcon",[["rect",{width:"14",height:"14",x:"8",y:"8",rx:"2",ry:"2",key:"17jyea"}],["path",{d:"M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2",key:"zix9uf"}]]),k=["title","aria-label"],f={key:2},C=i({__name:"CopyButton",props:{text:{},label:{default:"复制"},iconOnly:{type:Boolean,default:!1}},setup(r){const o=r,e=p(!1);async function c(){if(o.text)try{await navigator.clipboard.writeText(o.text),e.value=!0,setTimeout(()=>{e.value=!1},2e3)}catch{const t=document.createElement("textarea");t.value=o.text,document.body.appendChild(t),t.select(),document.execCommand("copy"),document.body.removeChild(t),e.value=!0,setTimeout(()=>{e.value=!1},2e3)}}return(t,v)=>(a(),l("button",{type:"button",onClick:m(c,["stop"]),title:e.value?"已复制到剪贴板":"点击复制","aria-label":e.value?"已复制":"复制",class:u(["inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-lg transition-all border select-none shadow-2xs active:scale-95",e.value?"bg-emerald-50 text-emerald-700 border-emerald-300 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-700":"bg-white text-slate-700 border-slate-200/90 hover:bg-slate-50 hover:text-brand-600 hover:border-brand-300 dark:bg-slate-800/90 dark:text-slate-300 dark:border-slate-700/80 dark:hover:bg-slate-700/90 dark:hover:text-brand-400 dark:hover:border-brand-600"])},[e.value?(a(),n(s(x),{key:0,class:"w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400"})):(a(),n(s(y),{key:1,class:"w-3.5 h-3.5 opacity-70"})),r.iconOnly?h("",!0):(a(),l("span",f,b(e.value?"已复制":r.label),1))],10,k))}});export{C as _};
