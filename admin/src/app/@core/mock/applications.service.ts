import { Injectable } from '@angular/core';
import { ApplicationsData, IApplication } from '../data/applications';
import { Observable, of } from 'rxjs';
import { delay } from 'rxjs/operators';

const mockApps: IApplication[] = [
];

@Injectable()
export class ApplicationsService extends ApplicationsData {
  getAll(): Observable<IApplication[]> {
    return of(mockApps).pipe(
      delay(100),
    );
  }
}
