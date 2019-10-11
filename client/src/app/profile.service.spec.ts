import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { HttpClient } from '@angular/common/http';

import { ProfileService } from './profile.service';
import { IProfile } from 'src/data/profile';

const mockProfile: IProfile = {
  intro: 'This is a mock profile!',
};

describe('ProfileService', () => {
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;
  let service: ProfileService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ HttpClientTestingModule ],
    });

    httpClient = TestBed.get(HttpClient);
    httpTestingController = TestBed.get(HttpTestingController);
    service = TestBed.get(ProfileService);
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it('is created', () => {
    expect(service).toBeTruthy();
  });

  it('gets profile data from the correct URL', () => {
    service.getProfile().subscribe(profile => {
      expect(profile).toEqual(mockProfile);
    });

    const req = httpTestingController.expectOne('api/profile');

    req.flush(mockProfile);
  });
});
