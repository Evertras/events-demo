import { TestBed, async, inject } from '@angular/core/testing';

import { AuthedGuard } from './authed.guard';

describe('AuthedGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AuthedGuard]
    });
  });

  it('should ...', inject([AuthedGuard], (guard: AuthedGuard) => {
    expect(guard).toBeTruthy();
  }));
});
