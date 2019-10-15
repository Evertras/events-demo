import { Injectable } from '@angular/core';
import { InMemoryDbService, RequestInfo, ResponseOptions } from 'angular-in-memory-web-api';
import { IProfile } from 'src/data/profile';
import { IOpenGame } from 'src/data/game';
import { LogService } from 'src/app/log.service';

const now = new Date();

@Injectable({
  providedIn: 'root'
})
export class InMemoryDataService implements InMemoryDbService {
  constructor(
    private log: LogService,
  ) { }

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

    /*
  get(_: RequestInfo) {
    this.log.debug('GET intercept');
  }

  post(reqInfo: RequestInfo) {
    this.log.debug('POST intercept');
    this.log.debug(JSON.stringify(reqInfo, null, 2));
  }

  put(_: RequestInfo) {
    this.log.debug('PUT intercept');
  }
     */

  responseInterceptor(resOptions: ResponseOptions, reqInfo: RequestInfo) {
    this.log.trace(`${reqInfo.method.toUpperCase()} ${reqInfo.req.url}`);

    return resOptions;
  }
}
