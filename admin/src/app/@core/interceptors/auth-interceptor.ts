import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent } from '@angular/common/http';
import { NbAuthService } from '@nebular/auth';
import { Observable } from 'rxjs';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  token: string;

  constructor(
    private authService: NbAuthService,
  ) {
    this.authService.onTokenChange().subscribe(token => {
      if (!token.isValid()) {
        this.token = '';
        return;
      }

      this.token = token.getValue();
    });
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    if (this.token) {
      const cloned = req.clone({
        headers: req.headers.set('X-Auth-Token', this.token),
      });

      return next.handle(cloned);
    } else {
      return next.handle(req);
    }
  }

}
