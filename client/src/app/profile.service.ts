import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { Observable, of } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { IProfile } from 'src/data/profile';
import { AuthService } from 'src/app/auth.service';
import { LogService } from './log.service';

const profileEndpoint = 'api/profile';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {

  constructor(
    private http: HttpClient,
    private auth: AuthService,
    private log: LogService,
  ) { }

  getProfile(): Observable<IProfile> {
    this.log.debug('Fetching profile');

    const headers = this.auth.authHeaders();

    const req = this.http.get<IProfile>(
      profileEndpoint,
      { headers },
    );

    return req.pipe(
      catchError(this.handleError<IProfile>('getProfile', {
        intro: '',
      }))
    );
  }

  private handleError<T>(method: string, safeRet: T) {
    return (error: any): Observable<T> => {
      this.log.error(`${method} failed: ${error.message || JSON.stringify(error)}`);

      return of(safeRet);
    };
  }
}
