/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue';
  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>;
  export default component;
}

declare module 'xlsx-js-style/dist/xlsx.bundle.js' {
  export type WorkBook = any;
  export type WorkSheet = any;
  export const utils: any;
  export const writeFile: any;
  export const write: any;
  const XLSX: any;
  export default XLSX;
}



