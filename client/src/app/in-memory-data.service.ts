import { Injectable } from '@angular/core';
import { InMemoryDbService } from 'angular-in-memory-web-api';
import { IProfile } from 'src/data/profile';
import { IOpenGame } from 'src/data/game';

const now = new Date();

@Injectable({
  providedIn: 'root'
})
export class InMemoryDataService implements InMemoryDbService {
  createDb() {
    const profile: IProfile = {
      intro: 'Hello!  This is an in-memory profile intro.',
    };

    const openGames: IOpenGame[] = [{
      id: 13,
      name: 'First in memory game',
      created: now,
    }];

    return {
      profile,
      openGames,
     };
  }

  constructor() { }
}
