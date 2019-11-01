import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { DevModule } from '../dev.module';

import { Header, HeaderData } from '../../../@core/data/headers';
import { HeadersComponent } from './headers.component';

const mockHeaders: Header[] = [
  {
    key: 'Mock-Header',
    value: 'Some mock value',
  },
  {
    key: 'X-Another-Header',
    value: 'abcedfasdkflja ljsadflk jsflkasjglkajsldjfhs df',
  },
];

const headerDataStub: Partial<HeaderData> = {
  getHeaders: () => of(mockHeaders),
};

describe('HeadersComponent', () => {
  let component: HeadersComponent;
  let fixture: ComponentFixture<HeadersComponent>;
  let nbListElement: HTMLElement;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        DevModule,
      ],
      providers: [
        {
          provide: HeaderData,
          useValue: headerDataStub,
        },
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HeadersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    nbListElement = fixture.nativeElement.querySelector('nb-list');
  });

  it('creates', () => {
    expect(component).toBeTruthy();
  });

  it('shows headers in a list', () => {
    const items = nbListElement.querySelectorAll('nb-list-item');
    const headerKeys = mockHeaders.map(h => `${h.key}: ${h.value}`);

    expect(items.length).toEqual(mockHeaders.length);
    items.forEach(el => {
      expect(headerKeys).toContain(el.textContent.trim());
    });
  });
});
