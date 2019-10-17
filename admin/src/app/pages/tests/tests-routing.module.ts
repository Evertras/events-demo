import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { TestsComponent } from './tests.component';
import { ListComponent } from './list/list.component';

const routes: Routes = [{
  path: '',
  component: TestsComponent,
  children: [
    {
      path: 'list',
      component: ListComponent,
    },
  ],
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class TestsRoutingModule { }

export const routedComponents = [
  ListComponent,
  TestsComponent,
];
