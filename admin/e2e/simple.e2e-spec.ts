import { browser } from 'protractor';

describe('simple', function() {
  it('loads the page', function() {
    browser.get('http://localhost:4200');

    expect(browser.getTitle()).toEqual('Admin');
  });
});

