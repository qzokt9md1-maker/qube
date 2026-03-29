// Fix for Node.js 22+: localStorage exists as a global but getItem is broken
// without --localstorage-file. This causes Next.js dev overlay to crash.
export async function register() {
  if (typeof window === "undefined" && typeof globalThis.localStorage !== "undefined") {
    // Running on server - delete the broken Node.js localStorage
    // @ts-expect-error - intentionally deleting broken global
    delete globalThis.localStorage;
  }
}
