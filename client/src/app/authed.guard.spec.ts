import { TestBed, inject } from '@angular/core/testing';

import { AuthedGuard } from './authed.guard';
import { AuthService } from 'src/app/auth.service';

let isAuthed = true;

const authServiceStub: Partial<AuthService> = {
  isAuthenticated() {
    return isAuthed;
  }
};

describe('AuthedGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        AuthedGuard,
        {
          provide: AuthService,
          useValue: authServiceStub,
        },
      ],
    });
  });

  it('checks the auth service for authentication', inject([AuthedGuard], (guard: AuthedGuard) => {
    isAuthed = true;

    expect(guard.canActivate({} as any, {} as any)).toBeTruthy();

    isAuthed = false;

    expect(guard.canActivate({} as any, {} as any)).toBeFalsy();
  }));
});
