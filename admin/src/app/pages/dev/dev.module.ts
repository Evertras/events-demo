import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { DevComponent } from './dev.component';
import { DevRoutingModule } from './dev-routing.module';

import { HeadersComponent } from './headers/headers.component';

@NgModule({
  imports: [
    CommonModule,
    DevRoutingModule,
  ],
  declarations: [
    DevComponent,
    HeadersComponent,
  ],
})
export class DevModule { }
