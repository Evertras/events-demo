import { Component, OnInit } from '@angular/core';
import { LogService } from 'src/app/log.service';
import {AuthService} from 'src/app/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  username: string;
  password: string;
  processing = false;

  constructor(
    private log: LogService,
    private auth: AuthService,
    private router: Router,
  ) { }

  ngOnInit() {
  }

  login() {
    if (this.username) {
      this.processing = true;
      this.log.debug('Logging in as ' + this.username);
      this.auth.login(this.username, this.password).subscribe(() => {
        this.processing = false;

        // TODO: Fancier redirects maybe?
        this.router.navigateByUrl('/');
      });
    }
  }

}
