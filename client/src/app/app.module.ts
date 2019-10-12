import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { HttpClientInMemoryWebApiModule } from 'angular-in-memory-web-api';

import { environment } from 'src/environments/environment';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ProfileComponent } from './profile/profile.component';
import { InMemoryDataService } from './in-memory-data.service';

const conditionalImports = [];

if (environment.inMemory) {
  conditionalImports.push(
    HttpClientInMemoryWebApiModule.forRoot(
      InMemoryDataService,
      {
        apiBase: 'api/',
        dataEncapsulation: false,
      })
  );
}

@NgModule({
  declarations: [
    AppComponent,
    ProfileComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    ...conditionalImports,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
