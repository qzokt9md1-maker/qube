const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/graphql";

class ApiClient {
  private accessToken: string | null = null;
  private refreshToken: string | null = null;

  constructor() {
    if (typeof window !== "undefined") {
      this.accessToken = localStorage.getItem("qube_token");
      this.refreshToken = localStorage.getItem("qube_refresh");
    }
  }

  setTokens(access: string, refresh: string) {
    this.accessToken = access;
    this.refreshToken = refresh;
    localStorage.setItem("qube_token", access);
    localStorage.setItem("qube_refresh", refresh);
  }

  clearTokens() {
    this.accessToken = null;
    this.refreshToken = null;
    localStorage.removeItem("qube_token");
    localStorage.removeItem("qube_refresh");
  }

  get isLoggedIn() {
    return !!this.accessToken;
  }

  async query<T = any>(operationName: string, variables?: Record<string, any>): Promise<T> {
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
