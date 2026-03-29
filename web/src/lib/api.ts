const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/graphql";

// Node.js 22+ has a broken global localStorage (no getItem method without --localstorage-file).
// We must check for both window AND working localStorage.
function isBrowserWithStorage(): boolean {
  try {
    return typeof window !== "undefined" && typeof window.document !== "undefined" && typeof window.localStorage?.getItem === "function";
  } catch {
    return false;
  }
}

function getStorage(key: string): string | null {
  if (!isBrowserWithStorage()) return null;
  try {
    return window.localStorage.getItem(key);
  } catch {
    return null;
  }
}

function setStorage(key: string, value: string) {
  if (!isBrowserWithStorage()) return;
  try {
    window.localStorage.setItem(key, value);
  } catch {}
}

function removeStorage(key: string) {
  if (!isBrowserWithStorage()) return;
  try {
    window.localStorage.removeItem(key);
  } catch {}
}

class ApiClient {
  private accessToken: string | null = null;
  private refreshToken: string | null = null;
  private initialized = false;

  private init() {
    if (this.initialized) return;
    this.initialized = true;
    this.accessToken = getStorage("qube_token");
    this.refreshToken = getStorage("qube_refresh");
  }

  setTokens(access: string, refresh: string) {
    this.accessToken = access;
    this.refreshToken = refresh;
    setStorage("qube_token", access);
    setStorage("qube_refresh", refresh);
  }

  clearTokens() {
    this.accessToken = null;
    this.refreshToken = null;
    removeStorage("qube_token");
    removeStorage("qube_refresh");
  }

  get isLoggedIn() {
    this.init();
    return !!this.accessToken;
  }

  async query<T = any>(operationName: string, variables?: Record<string, any>): Promise<T> {
    this.init();

    const res = await fetch(API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(this.accessToken ? { Authorization: `Bearer ${this.accessToken}` } : {}),
      },
      body: JSON.stringify({ operationName, query: "", variables: variables || {} }),
    });

    const json = await res.json();

    if (json.errors?.length > 0) {
      const msg = json.errors[0].message;
      if (msg === "unauthorized" && this.refreshToken) {
        const refreshed = await this.tryRefresh();
        if (refreshed) {
          return this.query(operationName, variables);
        }
      }
      throw new Error(msg);
    }

    return json.data as T;
  }

  private async tryRefresh(): Promise<boolean> {
    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          operationName: "RefreshToken",
          query: "",
          variables: { token: this.refreshToken },
        }),
      });
      const json = await res.json();
      if (json.data?.refreshToken) {
        const { accessToken, refreshToken } = json.data.refreshToken;
        this.setTokens(accessToken, refreshToken);
        return true;
      }
    } catch {}
    this.clearTokens();
    return false;
  }
}

export const api = new ApiClient();
