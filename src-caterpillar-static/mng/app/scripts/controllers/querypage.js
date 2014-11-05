'use strict';

/**
 * @ngdoc function
 * @name caterpillarApp.controller:QuerypageCtrl
 * @description
 * # QuerypageCtrl
 * Controller of the caterpillarApp
 */
angular.module('caterpillarApp')
  .controller('QuerypageCtrl', function ($scope, $window, Page, Rootpage) {
    $scope.pages = Page.query();
    $scope.criteria = {};

    $scope.makeRoot = function(page) {
      angular.forEach($scope.pages, function(another) {
        another._hover = null;
      });

      if (!$window.confirm('URLのルート(http://xxxx/)で表示されるページを変更します。\nよろしいですか？')) {
        return;
      }

      Rootpage.update({pageId: page.id})
      .$promise.then(function() {
        angular.forEach($scope.pages, function(another) {
          another.root = page === another ? true : false;
        });
      });
    };
  });
