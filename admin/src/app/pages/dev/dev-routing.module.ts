import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DevComponent } from './dev.component';
import { HeadersComponent } from './headers/headers.component';

const routes: Routes = [
  {
    path: '',
    component: DevComponent,
    children: [
      {
        path: 'headers',
        component: HeadersComponent,
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class DevRoutingModule {
}
