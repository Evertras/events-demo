import { Injectable } from '@angular/core';
import { InMemoryDbService, RequestInfo, ResponseOptions, getStatusText } from 'angular-in-memory-web-api';
import { IProfile } from 'src/data/profile';
import { IOpenGame } from 'src/data/game';
import { LogService } from 'src/app/log.service';
import { Observable, of } from 'rxjs';
import { HttpHeaders } from '@angular/common/http';

const now = new Date();
const fakeToken = 'abkjaug7dFxl';

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

  get(reqInfo: RequestInfo) {
    this.log.trace('GET ' + reqInfo.url);

    if (reqInfo.url !== 'api/auth/login' && !this.authHeaderValid(reqInfo)) {
      return reqInfo.utils.createResponse$(() => {
        const status = 401;
        const { headers, url } = reqInfo;

        const options: ResponseOptions = {
          headers,
          status,
          statusText: getStatusText(status),
          url,
        };

        return options;
      });
    }
  }

  private authHeaderValid(reqInfo: RequestInfo): boolean {
    // The typings here are not what they say they are...
    // https://github.com/angular/in-memory-web-api/issues/156
    const token = ((reqInfo.req as any).headers as HttpHeaders).get('X-Auth-Token');

    return token === fakeToken;
  }

  post(reqInfo: RequestInfo) {
    this.log.trace('POST ' + reqInfo.url);

    if (reqInfo.url === 'api/auth/login') {
      return reqInfo.utils.createResponse$(() => {
        const body = {token: fakeToken};
        const { headers, url } = reqInfo;
        const status = 200;

        const options: ResponseOptions = {
          body,
          headers,
          status,
          statusText: getStatusText(status),
          url,
        };

        return options;
      });
    }
  }

  put(reqInfo: RequestInfo) {
    this.log.trace('PUT ' + reqInfo.url);
  }

  // Leaving this declared here for future use, but right now it's essentially a pass-through
  responseInterceptor(resOptions: ResponseOptions, reqInfo: RequestInfo) {
    // this.log.trace(`${reqInfo.method.toUpperCase()} ${reqInfo.req.url}`);

    return resOptions;
  }
}
