import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { Observable, of, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';

import { LogService } from 'src/app/log.service';

const authHeaderKey = 'X-Auth-Token';
const authEndpoint = 'api/auth';
const loginEndpoint = `${authEndpoint}/login`;

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private tokenValue = '';

  constructor(
    private http: HttpClient,
    private log: LogService,
  ) { }

  isAuthenticated(): boolean {
    return !!this.tokenValue;
  }

  login(username: string, password: string): Observable<string> {
    this.log.debug('AuthService.login ' + loginEndpoint);

    const res = this.http.post<{ token: string }>(
      loginEndpoint,
      { username, password },
    );

    return res.pipe(
      map(r => r.token),
      tap(token => this.tokenValue = token),
      catchError(e => {
        this.log.warning(e);
        return throwError(e);
      }),
    );
  }

  authHeaders(existing?: HttpHeaders): HttpHeaders {
    if (existing) {
      return existing.set(authHeaderKey, this.tokenValue);
    }

    return new HttpHeaders({ [authHeaderKey]: this.tokenValue });
  }
}
