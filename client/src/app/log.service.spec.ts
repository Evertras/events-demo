import { TestBed } from '@angular/core/testing';

import { LogService, LogLevel } from './log.service';

describe('LogService', () => {
  let log: LogService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    log = TestBed.get(LogService);
  });

  it('is created', () => {
    expect(log).toBeTruthy();
  });

  it('sends log entries through to a subscriber with the correct levels', () => {
    const received: { [key in LogLevel]: number } = {
      [LogLevel.Trace]: 0,
      [LogLevel.Debug]: 0,
      [LogLevel.Info]: 0,
      [LogLevel.Warning]: 0,
      [LogLevel.Error]: 0,
    };

    const msgs: string[] = [];

    log.entries.subscribe(e => {
      received[e.level]++;
      msgs.push(e.msg);
    });

    expect(received[LogLevel.Trace]).toEqual(0);
    expect(received[LogLevel.Debug]).toEqual(0);
    expect(received[LogLevel.Info]).toEqual(0);
    expect(received[LogLevel.Warning]).toEqual(0);
    expect(received[LogLevel.Error]).toEqual(0);

    log.trace('Trace Test');
    log.debug('Debug Test');
    log.info('Info Test');
    log.warning('Warning Test');
    log.error('Error Test');

    expect(received[LogLevel.Trace]).toEqual(1);
    expect(received[LogLevel.Debug]).toEqual(1);
    expect(received[LogLevel.Info]).toEqual(1);
    expect(received[LogLevel.Warning]).toEqual(1);
    expect(received[LogLevel.Error]).toEqual(1);

    for (const msg of ['Trace', 'Debug', 'Info', 'Warning', 'Error'].map(s => `${s} Test`)) {
      expect(msgs).toContain(msg);
    }
  });

  it('sends log entries through to multiple subscribers with the correct levels', () => {
    const received: { [key in LogLevel]: number } = {
      [LogLevel.Trace]: 0,
      [LogLevel.Debug]: 0,
      [LogLevel.Info]: 0,
      [LogLevel.Warning]: 0,
      [LogLevel.Error]: 0,
    };

    log.entries.subscribe(e => {
      received[e.level]++;
    });

    log.entries.subscribe(e => {
      received[e.level]++;
    });

    log.trace('Trace');
    log.debug('Debug');
    log.info('Info');
    log.warning('Warning');
    log.error('Error');

    expect(received[LogLevel.Trace]).toEqual(2);
    expect(received[LogLevel.Debug]).toEqual(2);
    expect(received[LogLevel.Info]).toEqual(2);
    expect(received[LogLevel.Warning]).toEqual(2);
    expect(received[LogLevel.Error]).toEqual(2);
  });
});
