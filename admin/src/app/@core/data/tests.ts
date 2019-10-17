import {Observable} from 'rxjs';

export type VariableType = boolean | number | string;

export interface IVariable {
  name: string;
  type: 'bool' | 'int' | 'string';
  value: VariableType;
}

export interface IApplication {
  application_id: string;
  name: string;
  default_bundle_id: string;
  variables: IVariable[];
}

export interface ICohort {
  cohort_id: string;
  name: string;
  variables: IVariable[];
}

export enum TestStatus {
  available,
  active,
  finished,
  resolved,
}

export interface ITest {
  ab_test_id: string;
  application_id: string;
  name: string;
  maxUsers: number;
  status: TestStatus;
}

export abstract class TestsData {
  abstract getTests(): Observable<ITest[]>;
}
