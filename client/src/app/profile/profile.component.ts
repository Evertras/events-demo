import { Component, OnInit } from '@angular/core';
import { IProfile } from 'src/data/profile';
import { ProfileService } from 'src/app/profile.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {

  profile: IProfile = {
    intro: '',
  };

  constructor(
    private profileService: ProfileService,
  ) { }

  ngOnInit() {
    this.getProfile();
  }

  getProfile(): void {
    this.profileService.getProfile().subscribe(profile => this.profile = profile);
  }

}
