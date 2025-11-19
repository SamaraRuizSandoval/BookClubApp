const API_BASE = import.meta.env.VITE_API_BASE_URL; // or similar

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!res.ok) {
    // You can map errors here
    const text = await res.text();
    throw new Error(`HTTP ${res.status}: ${text}`);
  }

  // Handle 204 or empty bodies
  if (res.status === 204) return undefined as T;

  return res.json() as Promise<T>;
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
  // put, patch, delete, etc
};
