import { of as observableOf,  Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { Header, HeaderData } from '../data/headers';

@Injectable()
export class HeaderService extends HeaderData {

  private headers: Header[] = [
    {
      key: 'Accept',
      value: '*/*',
    },
    {
      key: 'X-User-Id',
      value: 'abc@def.com',
    },
    {
      key: 'User-Agent',
      value: 'In Memory Browser',
    },
    {
      key: 'X-Long-Header',
      value: 'This is a really long header and it has some interesting stuff in it I guess',
    },
  ];

  getHeaders(): Observable<Header[]> {
    return observableOf(this.headers);
  }
}
