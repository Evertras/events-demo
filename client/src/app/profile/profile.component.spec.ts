import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProfileComponent } from './profile.component';
import { ProfileService } from 'src/app/profile.service';
import { IProfile } from 'src/data/profile';
import { of } from 'rxjs';

const mockProfile: IProfile = {
  intro: 'This is a fake profile!',
};

let profileServiceStub: Partial<ProfileService>;

profileServiceStub = {
  getProfile: () => of(mockProfile),
};

describe('ProfileComponent', () => {
  let component: ProfileComponent;
  let fixture: ComponentFixture<ProfileComponent>;
  let compiled: any;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProfileComponent ],
      providers: [
        {
          provide: ProfileService,
          useValue: profileServiceStub,
        },
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProfileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    compiled = fixture.debugElement.nativeElement;
  });

  it('creates', () => {
    expect(component).toBeTruthy();
  });

  it('displays the intro', () => {
    expect(compiled.querySelector('.profile-intro').textContent).toContain(mockProfile.intro);
  });
});
