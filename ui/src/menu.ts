// Helper function to get user roles with fallback
function getUserRoles(): string | null {
  // First try to get from sessionStorage (primary)
  let roles = sessionStorage.getItem('roles');

  // If not found in sessionStorage, try to retrieve from localStorage backup
  if (!roles) {
    const backupRoles = localStorage.getItem('roles_backup');
    if (backupRoles) {
      // Restore from backup to sessionStorage
      sessionStorage.setItem('roles', backupRoles);
      roles = backupRoles;
    }
  }

  return roles;
}

export function isAdmin(): boolean {
  const s = getUserRoles();
  if (s) {
    return s.indexOf('Admin') > -1;
  }
  return false;
}

export function isLoggedIn(): boolean {
  const roles = getUserRoles();
  const isLoggedIn = roles != undefined && roles != null;
  return isLoggedIn;
}

export function isMaintainer(): boolean {
  const s = getUserRoles();
  if (s) {
    return s.indexOf('Maintainer') > -1;
  }
  return false;
}

export function isTester(): boolean {
  const s = getUserRoles();
  if (s) {
    return s.indexOf('Tester') > -1;
  }
  return false;
}
