import { Injectable } from '@angular/core';
import { TestsData, ITest, TestStatus } from '../data/tests';
import { Observable, of } from 'rxjs';
import { delay } from 'rxjs/operators';

const mockTests: ITest[] = [
  {
    ab_test_id: 'test1',
    application_id: 'app1',
    name: 'Test 1',
    maxUsers: 1000,
    status: TestStatus.available,
  },
  {
    ab_test_id: 'test2',
    application_id: 'app2',
    name: 'Another test for another app',
    maxUsers: 2000,
    status: TestStatus.available,
  },
  {
    ab_test_id: 'test3',
    application_id: 'app1',
    name: 'This one is active!',
    maxUsers: 5000,
    status: TestStatus.active,
  },
  {
    ab_test_id: 'test4',
    application_id: 'app1',
    name: 'This one finished...',
    maxUsers: 200,
    status: TestStatus.finished,
  },
];

@Injectable()
export class TestsService extends TestsData {
  getTests(): Observable<ITest[]> {
    return of(mockTests).pipe(
      delay(100),
    );
  }
}
