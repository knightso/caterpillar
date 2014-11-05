'use strict';

/**
 * @ngdoc function
 * @name caterpillarApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the caterpillarApp
 */
angular.module('caterpillarApp')
  .controller('MainCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
