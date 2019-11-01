import { Component, OnInit } from '@angular/core';

import { Header, HeaderData } from '../../../@core/data/headers';

@Component({
  selector: 'ngx-headers',
  templateUrl: './headers.component.html',
  styleUrls: ['./headers.component.scss'],
})
export class HeadersComponent implements OnInit {

  headers: Header[] = [];

  constructor(
    private headerService: HeaderData,
  ) {
    this.headerService.getHeaders().subscribe((headers) => this.headers = headers);
  }

  ngOnInit() {
  }

}
