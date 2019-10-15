import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';

import { HttpClientInMemoryWebApiModule } from 'angular-in-memory-web-api';

import { environment } from 'src/environments/environment';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ProfileComponent } from './profile/profile.component';
import { LoginComponent } from './login/login.component';
import { InMemoryDataService } from './in-memory-data.service';
import { HomeComponent } from './home/home.component';
import { DebugMessagesComponent } from './debug-messages/debug-messages.component';

const conditionalImports = [];

if (environment.inMemory) {
  conditionalImports.push(
    HttpClientInMemoryWebApiModule.forRoot(
      InMemoryDataService,
      {
        apiBase: 'api/',
        dataEncapsulation: false,
        delay: 200,
      }),
  );
}

@NgModule({
  declarations: [
    AppComponent,
    DebugMessagesComponent,
    HomeComponent,
    LoginComponent,
    ProfileComponent,
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    FormsModule,
    HttpClientModule,
    ...conditionalImports,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule { }
