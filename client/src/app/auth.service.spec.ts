import { TestBed } from '@angular/core/testing';

import { AuthService } from './auth.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import {
  HttpClientTestingModule,
  HttpTestingController
} from '@angular/common/http/testing';

describe('AuthService', () => {
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;
  let service: AuthService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ HttpClientTestingModule ],
    });

    httpClient = TestBed.get(HttpClient);
    httpTestingController = TestBed.get(HttpTestingController);
    service = TestBed.get(AuthService);
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it('is created', () => {
    expect(service).toBeTruthy();
  });

  function doLogin(): Promise<string> {
    const testUsername = 'mockuser';
    const testPassword = 'sekrits';
    const mockToken = 'abcdeftoken';

    return new Promise<string>((resolve) => {
      service.login(testUsername, testPassword).subscribe(token => {
        resolve(token);
      });

      const req = httpTestingController.expectOne(
        r => {
          return r.url === 'api/auth/login' &&
            r.body.username === testUsername &&
            r.body.password === testPassword;
        }
      );

      req.flush({ token: mockToken });
    });
  }

  it('logs in and creates a token', async () => {
    expect(service.isAuthenticated()).toBeFalsy();

    await doLogin();

    expect(service.isAuthenticated()).toBeTruthy();
  });

  it('creates a new HttpHeaders object when requested', async () => {
    await doLogin();

    const headers = service.authHeaders();

    expect(headers.get('X-Auth-Token')).toBeTruthy();
  });

  it('adds to an existing HttpHeaders object when requested', async () => {
    await doLogin();

    const existingHeaders = new HttpHeaders({ 'X-SomeOtherHeader': 'something' });
    const headers = service.authHeaders(existingHeaders);

    expect(headers.get('X-Auth-Token')).toBeTruthy();
    expect(headers.get('X-SomeOtherHeader')).toEqual('something');
  });
});
