export const apiBase = process.env.NEXT_PUBLIC_API_URL || "/api";

function getAuthHeader() {
  if (typeof window === "undefined") return {};
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function apiFetch<T = unknown>(path: string, init: RequestInit = {}): Promise<T> {
  const url = path.startsWith("http") ? path : `${apiBase}${path}`;
  const headers = new Headers(init.headers as HeadersInit | undefined);
  headers.set("Content-Type", "application/json");
  const auth = getAuthHeader() as Record<string, string>;
  Object.entries(auth).forEach(([k, v]) => headers.set(k, v));

  const res = await fetch(url, { ...init, headers });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${res.status} ${res.statusText}: ${text}`);
  }
  try {
    return (await res.json()) as T;
  } catch {
    return undefined as T;
  }
}
