import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardComponent } from './dashboard.component';
import { NbCardModule } from '@nebular/theme';

// import { Subject } from 'rxjs';
// import { CompileShallowModuleMetadata } from '@angular/compiler';

/*
const mockEntries: Subject<ILogEntry> = new Subject<ILogEntry>();

const logServiceStub: Partial<LogService> = {
  trace: () => {},
  entries: mockEntries,
};
*/

describe('DashboardComponent', () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;
  let compiled: HTMLElement;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DashboardComponent ],
      imports: [ NbCardModule ],
      providers: [
        /*
        {
          provide: LogService,
          useValue: logServiceStub,
        }
         */
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    compiled = fixture.debugElement.nativeElement;
  });

  it('creates', () => {
    expect(component).toBeTruthy();
  });

  /*
  it('shows the table when there are some messages', () => {
    mockEntries.next({
      level: LogLevel.Info,
      msg: 'This is a message',
    });
    expect(component.entries.length).toEqual(1);
    fixture.detectChanges();
    expect(compiled.querySelector('table')).toBeTruthy();
  });

  it('hides the table when messages are cleared', () => {
    mockEntries.next({
      level: LogLevel.Info,
      msg: 'This is a message',
    });
    expect(component.entries.length).toEqual(1, 'Did not receive initial message');
    fixture.detectChanges();
    expect(compiled.querySelector('table')).toBeTruthy('Should show table now that messages exist');

    const button: HTMLButtonElement = compiled.querySelector('button.clear');

    expect(button).toBeDefined();

    button.click();

    expect(component.entries.length).toEqual(0, 'Did not clear messages');
    fixture.detectChanges();
    expect(compiled.querySelector('table')).toBeFalsy('Should not show table after clear');
  });
  */
});
