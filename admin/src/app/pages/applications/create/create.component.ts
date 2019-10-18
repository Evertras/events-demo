import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { IVariable } from '../../../@core/data/applications';

@Component({
  selector: 'ngx-create',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.scss'],
})
export class CreateComponent implements OnInit {

  infoForm: FormGroup;
  variables: IVariable[] = [];

  constructor(
    private fb: FormBuilder,
  ) {
    this.infoForm = this.fb.group({
      nameCtrl: ['', Validators.required],
    });
  }

  ngOnInit() {
  }

  onInfoSubmit() {
    this.infoForm.markAsDirty();
  }

  varsChanged(variables: IVariable[]) {
    this.variables = variables;
  }
}
