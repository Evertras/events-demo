import {Observable} from 'rxjs';

export type VariableType = boolean | number | string;

export interface IVariableOverride {
  name: string;
  type: 'bool' | 'int' | 'string';
  value: VariableType;
}

export interface ICohort {
  cohort_id: string;
  name: string;
  variables: IVariableOverride[];
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
