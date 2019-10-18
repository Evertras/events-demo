import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { LocalDataSource } from 'ng2-smart-table';
import { IVariable } from '../../../@core/data/applications';

@Component({
  selector: 'ngx-variables',
  templateUrl: './variables.component.html',
  styleUrls: ['./variables.component.scss'],
})
export class VariablesComponent implements OnInit {

  @Output() variablesChanged = new EventEmitter<IVariable[]>();

  addForm: FormGroup;

  tableSettings = {
    hideSubHeader: true,
    add: {
      addButtonContent: '<i class="nb-plus"></i>',
      createButtonContent: '<i class="nb-checkmark"></i>',
      cancelButtonContent: '<i class="nb-close"></i>',
    },
    edit: {
      editButtonContent: '<i class="nb-edit"></i>',
      saveButtonContent: '<i class="nb-checkmark"></i>',
      cancelButtonContent: '<i class="nb-close"></i>',
    },
    delete: {
      deleteButtonContent: '<i class="nb-trash"></i>',
    },
    columns: {
      name: {
        title: 'Name',
        type: 'string',
      },
      type: {
        title: 'Type',
        type: 'string',
        editor: {
          type: 'list',
          config: {
            list: ['string', 'int', 'bool'].map(t => ({ title: t, value: t })),
          },
        },
      },
      value: {
        title: 'Default',
        type: 'string',
      },
    },
    noDataMessage: 'No variables added yet',
  };

  tableSource: LocalDataSource = new LocalDataSource();

  constructor(
    private fb: FormBuilder,
  ) {
    this.addForm = this.fb.group({
      name: ['', Validators.required],
      default: ['', Validators.required],
      type: 'string',
    });

    this.tableSource.onChanged().subscribe(
      async () => this.variablesChanged.emit(await this.tableSource.getAll()),
    );
  }

  ngOnInit() {
  }

  async add() {
    this.tableSource.append({
      name: this.addForm.value.name,
      value: this.addForm.value.default,
      type: this.addForm.value.type,
    });

    this.variablesChanged.emit(await this.tableSource.getAll());

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
