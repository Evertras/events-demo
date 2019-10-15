import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { map, catchError, tap } from 'rxjs/operators';
import { Observable, of, throwError } from 'rxjs';

import { LogService } from 'src/app/log.service';

const authHeaderKey = 'X-Auth-Token';
const authEndpoint = 'api/auth';
const loginEndpoint = `${authEndpoint}/login`;

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(
    private http: HttpClient,
    private log: LogService,
  ) { }

  isAuthenticated(): boolean {
    // TODO: actual auth
    return true;
  }

  login(username: string, password: string): Observable<string> {
    this.log.debug('AuthService.login ' + loginEndpoint);

    console.log(this.http);

    const res = this.http.post<{ token: string }>(
      loginEndpoint,
      { username, password },
    );

    return res.pipe(
      map(r => r.token),
      tap(_ => this.log.trace('hi')),
      catchError(e => {
        this.log.warning(e);
        return throwError(e);
      }),
    );
  }

  authHeaders(existing?: HttpHeaders): HttpHeaders {
    // TODO: actual auth
    if (existing) {
      return existing.set(authHeaderKey, 'totally-secure');
    }

    return new HttpHeaders({ [authHeaderKey]: 'totally-secure' });
  }

  private handleError<T>(method: string, safeRet: T) {
    return (error: any): Observable<T> => {
      this.log.error(`${method} failed: ${error.message || error}`);

      return of(safeRet);
    };
  }
}
