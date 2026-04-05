import type { User } from "@/entities/user/model";

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
}

const TOKENS_KEY = "happyhouse.tokens";
const USER_KEY = "happyhouse.user";
const HOUSE_KEY = "happyhouse.houseId";

export function loadTokens(): AuthTokens | null {
  const raw = localStorage.getItem(TOKENS_KEY);
  return raw ? (JSON.parse(raw) as AuthTokens) : null;
}

export function saveTokens(tokens: AuthTokens) {
  localStorage.setItem(TOKENS_KEY, JSON.stringify(tokens));
}

export function clearTokens() {
  localStorage.removeItem(TOKENS_KEY);
}

export function loadUser(): User | null {
  const raw = localStorage.getItem(USER_KEY);
  return raw ? (JSON.parse(raw) as User) : null;
}

export function saveUser(user: User) {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearUser() {
  localStorage.removeItem(USER_KEY);
}

export function loadSelectedHouseId(): number | null {
  const raw = localStorage.getItem(HOUSE_KEY);
  return raw ? Number(raw) : null;
}

export function saveSelectedHouseId(houseId: number) {
  localStorage.setItem(HOUSE_KEY, String(houseId));
}

export function clearSelectedHouseId() {
  localStorage.removeItem(HOUSE_KEY);
}

export function clearSession() {
  clearTokens();
  clearUser();
  clearSelectedHouseId();
}
