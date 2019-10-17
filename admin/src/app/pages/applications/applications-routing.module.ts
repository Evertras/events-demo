import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ApplicationsComponent } from './applications.component';
import { ListComponent } from './list/list.component';
import { CreateComponent } from './create/create.component';

const routes: Routes = [{
  path: '',
  component: ApplicationsComponent,
  children: [
    {
      path: 'list',
      component: ListComponent,
    },
    {
      path: 'create',
      component: CreateComponent,
    },
  ],
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ApplicationsRoutingModule { }

export const routedComponents = [
  CreateComponent,
  ListComponent,
  ApplicationsComponent,
];
