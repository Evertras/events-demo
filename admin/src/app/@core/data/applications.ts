import { Observable } from 'rxjs';

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

export abstract class ApplicationsData {
  abstract getAll(): Observable<IApplication[]>;
}
