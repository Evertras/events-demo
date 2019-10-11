import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { IProfile } from 'src/data/profile';
import { Observable } from 'rxjs';

const profileEndpoint = 'api/profile';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {

  constructor(
    private http: HttpClient,
  ) { }

  getProfile(): Observable<IProfile> {
    return this.http.get<IProfile>(profileEndpoint);
  }
}
