import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

export enum LogLevel {
  Trace = 0,
  Debug = 1,
  Info = 2,
  Warning = 3,
  Error = 4,
}

export interface ILogEntry {
  level: LogLevel;
  msg: string;
}

@Injectable({
  providedIn: 'root'
})
export class LogService {

  entries: Subject<ILogEntry> = new Subject<ILogEntry>();

  constructor() { }

  trace(msg: string) {
    return this.log(LogLevel.Trace, msg);
  }

  debug(msg: string) {
    return this.log(LogLevel.Debug, msg);
  }

  info(msg: string) {
    return this.log(LogLevel.Info, msg);
  }

  warning(msg: string) {
    return this.log(LogLevel.Warning, msg);
  }

  error(msg: string) {
    return this.log(LogLevel.Error, msg);
  }

  log(level: LogLevel, msg: string) {
    this.entries.next({
      level,
      msg,
    });
  }
}
