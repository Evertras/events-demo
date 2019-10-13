import { Component, OnInit } from '@angular/core';
import { LogService, ILogEntry, LogLevel } from '../log.service';

const levelNames: {[key in LogLevel]: string} = {
  [LogLevel.Trace]: 'Trace',
  [LogLevel.Debug]: 'Debug',
  [LogLevel.Info]: 'Info',
  [LogLevel.Warning]: 'Warning',
  [LogLevel.Error]: 'Error',
};

@Component({
  selector: 'app-debug-messages',
  templateUrl: './debug-messages.component.html',
  styleUrls: ['./debug-messages.component.scss']
})
export class DebugMessagesComponent implements OnInit {

  entries: ILogEntry[] = [];

  constructor(
    private log: LogService,
  ) { }

  ngOnInit() {
    this.log.entries.subscribe(e => this.entries.push(e));
    this.log.trace('Initialized debug-messages component');
  }

  getLevelName(level: LogLevel) {
    return levelNames[level];
  }

  clearEntries() {
    this.entries = [];
  }

}
