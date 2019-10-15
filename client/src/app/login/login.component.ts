import { Component, OnInit } from '@angular/core';
import { LogService } from 'src/app/log.service';
import {AuthService} from 'src/app/auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  username: string;
  password: string;

  constructor(
    private log: LogService,
    private auth: AuthService,
  ) { }

  ngOnInit() {
  }

  login() {
    if (this.username) {
      this.log.debug('Logging in as ' + this.username);
      this.auth.login(this.username, this.password).subscribe();
    }
  }

}
