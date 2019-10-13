import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs';

import { IProfile } from 'src/data/profile';
import { AuthService } from 'src/app/auth.service';

const profileEndpoint = 'api/profile';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {

  constructor(
    private http: HttpClient,
    private auth: AuthService,
  ) { }

  getProfile(): Observable<IProfile> {
    return this.http.get<IProfile>(profileEndpoint, { headers: this.auth.authHeaders() });
  }
}
