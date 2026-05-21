export interface FeatureFlagRecord {
  id: string;
  app_id: string;
  name: string;
  description: string;
  enabled: boolean;
  rollout_percentage: number;
}

interface CacheEntry {
  value: boolean;
  expiresAt: number;
}

const CACHE_TTL_MS = 5000;

/** URL por defecto: servicio local con Docker Compose. */
export const DEFAULT_FLAG_SERVICE_URL = "http://localhost:8080";

/** FNV-1a 32-bit — mismo algoritmo que el servidor Go (hash/fnv). */
function fnv1a32(input: string): number {
  let hash = 0x811c9dc5;
  for (let i = 0; i < input.length; i++) {
    hash ^= input.charCodeAt(i);
    hash = Math.imul(hash, 0x01000193);
  }
  return hash >>> 0;
}

function evaluateRollout(
  flagName: string,
  userId: string,
  rolloutPercentage: number
): boolean {
  if (rolloutPercentage >= 100) return true;
  if (rolloutPercentage <= 0) return false;

  const key = `${flagName}:${userId || "anonymous"}`;
  const bucket = fnv1a32(key) % 100;
  return bucket < rolloutPercentage;
}

export class FeatureFlagClient {
  private baseURL: string;
  private appId: string;
  private cache = new Map<string, CacheEntry>();

  constructor(
    appId: string,
    baseURL: string = DEFAULT_FLAG_SERVICE_URL
  ) {
    this.baseURL = baseURL.replace(/\/$/, "");
    this.appId = appId;
  }

  async isEnabled(flagName: string, userId?: string): Promise<boolean> {
    const cacheKey = `${flagName}:${userId ?? ""}`;
    const cached = this.cache.get(cacheKey);
    if (cached && cached.expiresAt > Date.now()) {
      return cached.value;
    }

    const headers: Record<string, string> = {};
    if (userId) {
      headers["X-User-Id"] = userId;
    }

    const res = await fetch(
      `${this.baseURL}/api/v1/public/flag/${this.appId}/${encodeURIComponent(flagName)}`,
      { headers }
    );

    if (!res.ok) {
      throw new Error(`Feature flag request failed: ${res.status}`);
    }

    const data = (await res.json()) as { enabled: boolean };
    this.cache.set(cacheKey, {
      value: data.enabled,
      expiresAt: Date.now() + CACHE_TTL_MS,
    });
    return data.enabled;
  }

  async getAllFlags(userId?: string): Promise<Record<string, boolean>> {
    const res = await fetch(
      `${this.baseURL}/api/v1/apps/${this.appId}/flags`
    );
    if (!res.ok) {
      throw new Error(`Failed to list flags: ${res.status}`);
    }

    const flags = (await res.json()) as FeatureFlagRecord[];
    const result: Record<string, boolean> = {};

    for (const flag of flags) {
      if (!flag.enabled) {
        result[flag.name] = false;
        continue;
      }
      result[flag.name] = evaluateRollout(
        flag.name,
        userId ?? "",
        flag.rollout_percentage
      );
    }
    return result;
  }
}

/**
 * Ejemplo de uso (FlagFlow local con Docker):
 *
 * ```ts
 * const client = new FeatureFlagClient("550e8400-e29b-41d4-a716-446655440000");
 * if (await client.isEnabled("dark-mode", "user-123")) {
 *   console.log("Dark mode activo");
 * }
 * ```
 */
