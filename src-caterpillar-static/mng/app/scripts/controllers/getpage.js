'use strict';

/**
 * @ngdoc function
 * @name caterpillarApp.controller:GetpageCtrl
 * @description
 * # GetpageCtrl
 * Controller of the caterpillarApp
 */
angular.module('caterpillarApp')
  .controller('GetpageCtrl', function ($scope, $routeParams, Page) {
    $scope.page = Page.get({
      id : $routeParams.id ? $routeParams.id : null
    });
  });
