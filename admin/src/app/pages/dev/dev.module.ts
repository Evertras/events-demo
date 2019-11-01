import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NbCardModule, NbListModule } from '@nebular/theme';

import { DevComponent } from './dev.component';
import { DevRoutingModule } from './dev-routing.module';

import { HeadersComponent } from './headers/headers.component';

@NgModule({
  imports: [
    CommonModule,
    DevRoutingModule,
    NbCardModule,
    NbListModule,
  ],
  declarations: [
    DevComponent,
    HeadersComponent,
  ],
})
export class DevModule { }
