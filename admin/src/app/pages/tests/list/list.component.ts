import { LocalDataSource } from 'ng2-smart-table';
import { Component, OnInit } from '@angular/core';

import { TestsData } from '../../../@core/data/tests';

@Component({
  selector: 'ngx-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ListComponent implements OnInit {

  settings = {
    actions: false,
    columns: {
      ab_test_id: {
        title: 'ID',
        type: 'string',
      },
      name: {
        title: 'Name',
        type: 'string',
      },
      maxUsers: {
        title: 'Max Users',
        type: 'number',
      },
    },
  };

  source: LocalDataSource = new LocalDataSource();

  constructor(
    private testsService: TestsData,
  ) {
  }

  ngOnInit() {
    this.testsService.getTests().subscribe(tests => {
      this.source.load(tests);
    });
  }

}
