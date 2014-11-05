'use strict';

describe('Service: Rootpage', function () {

  // load the service's module
  beforeEach(module('caterpillarApp'));

  // instantiate service
  var Rootpage;
  beforeEach(inject(function (_Rootpage_) {
    Rootpage = _Rootpage_;
  }));

  it('should do something', function () {
    expect(!!Rootpage).toBe(true);
  });

});
