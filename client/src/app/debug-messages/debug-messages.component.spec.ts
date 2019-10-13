import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DebugMessagesComponent } from './debug-messages.component';
import { LogService, ILogEntry, LogLevel } from '../log.service';
import { Subject } from 'rxjs';
import { CompileShallowModuleMetadata } from '@angular/compiler';

const mockEntries: Subject<ILogEntry> = new Subject<ILogEntry>();

const logServiceStub: Partial<LogService> = {
  trace: () => {},
  entries: mockEntries,
};

describe('DebugMessagesComponent', () => {
  let component: DebugMessagesComponent;
  let fixture: ComponentFixture<DebugMessagesComponent>;
  let compiled: HTMLElement;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DebugMessagesComponent ],
      providers: [
        {
          provide: LogService,
          useValue: logServiceStub,
        }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DebugMessagesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    compiled = fixture.debugElement.nativeElement;
  });

  it('creates', () => {
    expect(component).toBeTruthy();
  });

  it('initially hides the table when there are no messages to begin with', () => {
    expect(compiled.querySelector('table')).toBeFalsy();
  });

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
});
