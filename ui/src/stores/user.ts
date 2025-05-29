import { reactive } from 'vue';
import http from '@/common-http';

interface UserState {
  roles: string[];
  userId: number | null;
  email: string | null;
}

export const userState = reactive<UserState>({
  roles: [],
  userId: null,
  email: null
});

export function setUser(userId: number, email: string, roles: string) {
  userState.userId = userId;
  userState.email = email;
  userState.roles = roles ? roles.split(',') : [];
}

export function clearUser() {
  userState.roles = [];
  userState.userId = null;
  userState.email = null;
}

export async function fetchCurrentUser() {
  try {
    const response = await http.get('/api/v1/auth/me');
    if (response.data && response.data.data) {
      const data = response.data.data;
      setUser(data.userId, data.email, data.roles);
    }
  } catch {
    clearUser();
  }
}

export function isLoggedIn(): boolean {
  return userState.roles.length > 0;
}
export function isAdmin(): boolean {
  return userState.roles.includes('Admin');
}
export function isMaintainer(): boolean {
  return userState.roles.includes('Maintainer');
}
export function isTester(): boolean {
  return userState.roles.includes('Tester');
}
