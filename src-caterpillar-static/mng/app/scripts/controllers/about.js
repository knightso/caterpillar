'use strict';

/**
 * @ngdoc function
 * @name caterpillarApp.controller:AboutCtrl
 * @description
 * # AboutCtrl
 * Controller of the caterpillarApp
 */
angular.module('caterpillarApp')
  .controller('AboutCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
