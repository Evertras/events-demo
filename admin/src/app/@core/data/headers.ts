import { Observable } from 'rxjs';

export interface Header {
  key: string;
  value: string;
}

export abstract class HeaderData {
  abstract getHeaders(): Observable<Header[]>;
}

