import { Injectable } from '@angular/core';
import {InMemoryDbService} from 'angular-in-memory-web-api';
import {IProfile} from 'src/data/profile';

@Injectable({
  providedIn: 'root'
})
export class InMemoryDataService implements InMemoryDbService {
  createDb() {
    const profile: IProfile = {
      intro: 'Hello!',
    };

    return { profile };
  }

  constructor() { }
}
