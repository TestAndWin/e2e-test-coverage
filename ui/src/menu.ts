import { userState } from './stores/user';

export function isAdmin(): boolean {
  return userState.roles.includes('Admin');
}

export function isLoggedIn(): boolean {
  return userState.roles.length > 0;
}

export function isMaintainer(): boolean {
  return userState.roles.includes('Maintainer');
}

export function isTester(): boolean {
  return userState.roles.includes('Tester');
}
