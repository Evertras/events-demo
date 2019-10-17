import { Component, OnInit, Input } from '@angular/core';
import { IVariable } from '../../../@core/data/applications';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

@Component({
  selector: 'ngx-variables',
  templateUrl: './variables.component.html',
  styleUrls: ['./variables.component.scss'],
})
export class VariablesComponent implements OnInit {

  @Input() variables: IVariable[];

  addForm: FormGroup;

  constructor(
    private fb: FormBuilder,
  ) {
    this.addForm = this.fb.group({
      name: ['', Validators.required],
      default: ['', Validators.required],
      type: 'string',
    });
  }

  ngOnInit() {
  }

  add() {
    this.addForm.reset({
      type: 'string',
    });
  }

  resetDefault() {
    let val: string = '';

    if (this.addForm.value.type === 'bool') {
      val = 'true';
    }

    this.addForm.setValue(Object.assign(this.addForm.value, {
      default: val,
    }));
  }

}
