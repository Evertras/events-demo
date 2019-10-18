import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { Ng2SmartTableModule } from 'ng2-smart-table';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ApplicationsRoutingModule, routedComponents } from './applications-routing.module';
import {
  NbButtonModule,
  NbCardModule,
  NbListModule,
  NbSelectModule,
  NbStepperModule,
} from '@nebular/theme';
import { VariablesComponent } from './variables/variables.component';
import { ThemeModule } from '../../@theme/theme.module';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    NbButtonModule,
    NbCardModule,
    NbListModule,
    NbSelectModule,
    NbStepperModule,
    Ng2SmartTableModule,
    ApplicationsRoutingModule,
    ThemeModule,
  ],
  declarations: [
    ...routedComponents,
    VariablesComponent,
  ],
})
export class ApplicationsModule { }
