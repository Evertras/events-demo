import { ModuleWithProviders, NgModule, Optional, SkipSelf } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NbSecurityModule, NbRoleProvider } from '@nebular/security';
import { of as observableOf } from 'rxjs';

import { httpInterceptorProviders } from './interceptors';
import { throwIfAlreadyLoaded } from './module-import-guard';
import { AnalyticsService } from './utils';
import { UserData } from './data/users';
import { HeaderData } from './data/headers';
import { UserService } from './mock/users.service';
import { HeaderService } from './mock/headers.service';
import { MockDataModule } from './mock/mock-data.module';

import {
  NbAuthModule,
  NbAuthJWTToken,
  NbDummyAuthStrategy,
  NbPasswordAuthStrategy,
} from '@nebular/auth';

const socialLinks = [
  /*
  {
    url: 'https://github.com/akveo/nebular',
    target: '_blank',
    icon: 'github',
  },
   */
];

const DATA_SERVICES = [
  { provide: UserData, useClass: UserService },
  { provide: HeaderData, useClass: HeaderService },
];

const authStrategy = false ?
  NbDummyAuthStrategy.setup({
    name: 'email',
    delay: 1000,
  }) :
  NbPasswordAuthStrategy.setup({
    name: 'email',
    login: {
      method: 'POST',
    },
    token: {
      class: NbAuthJWTToken,
      key: 'token',
    },
  });

export class NbSimpleRoleProvider extends NbRoleProvider {
  getRole() {
    // here you could provide any role based on any auth flow
    return observableOf('guest');
  }
}

export const NB_CORE_PROVIDERS = [
  ...MockDataModule.forRoot().providers,
  ...DATA_SERVICES,
  ...httpInterceptorProviders,

  ...NbAuthModule.forRoot({
    strategies: [
      authStrategy,
    ],
    forms: {
      login: {
        rememberMe: false,
      },
      register: {
        socialLinks: socialLinks,
      },
    },
  }).providers,

  NbSecurityModule.forRoot({
    accessControl: {
      guest: {
        view: '*',
      },
      user: {
        parent: 'guest',
        create: '*',
        edit: '*',
        remove: '*',
      },
    },
  }).providers,

  {
    provide: NbRoleProvider,
    useClass: NbSimpleRoleProvider,
  },

  AnalyticsService,
];

@NgModule({
  imports: [
    CommonModule,
  ],
  exports: [
    NbAuthModule,
  ],
  declarations: [],
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    throwIfAlreadyLoaded(parentModule, 'CoreModule');
  }

  static forRoot(): ModuleWithProviders {
    return <ModuleWithProviders>{
      ngModule: CoreModule,
      providers: [
        ...NB_CORE_PROVIDERS,
      ],
    };
  }
}
