import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController
} from '@angular/common/http/testing';
import { HttpClient } from '@angular/common/http';

import { GameService } from './game.service';
import { IOpenGame } from 'src/data/game';

const now = new Date();

const mockOpenGames: IOpenGame[] = [
  {
    id: 17,
    name: 'An open game',
    created: new Date(now.getTime() - 60 * 1000),
  },
  {
    id: 481,
    name: 'テスト',
    created: new Date(now.getTime() - 120 * 1000),
  },
];

describe('GameService', () => {
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;
  let service: GameService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ HttpClientTestingModule ],
    });

    httpClient = TestBed.get(HttpClient);
    httpTestingController = TestBed.get(HttpTestingController);
    service = TestBed.get(GameService);
  });

  it('is created', () => {
    expect(service).toBeTruthy();
  });

  it('gets an open game list', () => {
    service.getOpenGames().subscribe(games => {
      expect(games.length).toBe(mockOpenGames.length);
    });

    const req = httpTestingController.expectOne(
      r => r.url === 'api/games/open'
        && r.headers.has('X-Auth-Token')
    );

    req.flush(mockOpenGames);
  });
});
