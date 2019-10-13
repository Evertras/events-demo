import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';

import { AuthService } from './auth.service';
import { IOpenGame } from 'src/data/game';

const openGameListEndpoint = 'api/games/open';

@Injectable({
  providedIn: 'root'
})
export class GameService {

  constructor(
    private http: HttpClient,
    private auth: AuthService,
  ) { }

  // No paging or filtering for now for simplicity
  getOpenGames() {
    return this.http.get<IOpenGame[]>(openGameListEndpoint, {
      headers: this.auth.authHeaders(),
    });
  }
}
