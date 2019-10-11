import { Injectable } from '@angular/core';
import { HttpHeaders } from '@angular/common/http';

const authHeaderKey = 'X-Auth-Token';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  isAuthenticated(): boolean {
    // TODO: actual auth
    return true;
  }

  authHeaders(existing?: HttpHeaders): HttpHeaders {
    // TODO: actual auth
    if (existing) {
      return existing.set(authHeaderKey, 'totally-secure');
    }

    return new HttpHeaders({ [authHeaderKey]: 'totally-secure' });
  }

  constructor() { }
}
