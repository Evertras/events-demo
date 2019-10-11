import { TestBed } from '@angular/core/testing';

import { AuthService } from './auth.service';
import { HttpHeaders } from '@angular/common/http';

describe('AuthService', () => {
  let service: AuthService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.get(AuthService);
  });

  it('is created', () => {
    expect(service).toBeTruthy();
  });

  // TODO: actual auth
  it('creates a new HttpHeaders object when requested', () => {
    const headers = service.authHeaders();

    expect(headers.get('X-Auth-Token')).toBeTruthy();
  });

  // TODO: actual auth
  it('adds to an existing HttpHeaders object when requested', () => {
    const existingHeaders = new HttpHeaders({ 'X-SomeOtherHeader': 'something' });
    const headers = service.authHeaders(existingHeaders);

    expect(headers.get('X-Auth-Token')).toBeTruthy();
    expect(headers.get('X-SomeOtherHeader')).toEqual('something');
  });
});
