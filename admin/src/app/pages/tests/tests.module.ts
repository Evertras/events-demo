import { NgModule } from '@angular/core';
import { Ng2SmartTableModule } from 'ng2-smart-table';
import { ListComponent } from './list/list.component';
import { TestsRoutingModule, routedComponents } from './tests-routing.module';
import { NbCardModule } from '@nebular/theme';

@NgModule({
  imports: [
    NbCardModule,
    Ng2SmartTableModule,
    TestsRoutingModule,
  ],
  declarations: [
    ...routedComponents,
    ListComponent,
  ],
})
export class TestsModule { }
