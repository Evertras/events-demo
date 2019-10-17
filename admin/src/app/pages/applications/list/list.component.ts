import { LocalDataSource } from 'ng2-smart-table';
import { Component, OnInit } from '@angular/core';

import { ApplicationsData } from '../../../@core/data/applications';

@Component({
  selector: 'ngx-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ListComponent implements OnInit {

  settings = {
    actions: false,
    columns: {
      application_id: {
        title: 'ID',
        type: 'string',
      },
      name: {
        title: 'Name',
        type: 'string',
      },
    },
  };

  source: LocalDataSource = new LocalDataSource();

  constructor(
    private applicationsService: ApplicationsData,
  ) {
  }

  ngOnInit() {
    this.applicationsService.getAll().subscribe(apps => {
      this.source.load(apps);
    });
  }

}
