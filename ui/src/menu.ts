export function isAdmin(): boolean {
  const s = sessionStorage.getItem('roles');
  if (s) {
    return s.indexOf('Admin') > -1;
  }
  return false;
}

export function isLoggedIn(): boolean {
  return sessionStorage.getItem('roles') != undefined;
}

export function isMaintainer(): boolean {
  const s = sessionStorage.getItem('roles');
  if (s) {
    return s.indexOf('Maintainer') > -1;
  }
  return false;
}

export function isTester(): boolean {
  const s = sessionStorage.getItem('roles');
  if (s) {
    return s.indexOf('Tester') > -1;
  }
  return false;
}
