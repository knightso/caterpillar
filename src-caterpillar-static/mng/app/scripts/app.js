'use strict';

/**
 * @ngdoc overview
 * @name caterpillarApp
 * @description
 * # caterpillarApp
 *
 * Main module of the application.
 */
angular
  .module('caterpillarApp', [
    'ngCookies',
    'ngResource',
    'ui.bootstrap',
    'ngRoute'
  ])
  .config(function ($routeProvider, $httpProvider) {
    $routeProvider
      .when('/', {
        redirectTo: '/queryPage'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
      })
      .when('/putPage/:id', {
        templateUrl: 'views/putpage.html',
        controller: 'PutpageCtrl',
        resolve: {method: function(){return 'PUT';}}
      })
      .when('/postPage', {
        templateUrl: 'views/putpage.html',
        controller: 'PutpageCtrl',
        resolve: {method: function(){return 'POST';}}
      })
      .when('/queryPage', {
        templateUrl: 'views/querypage.html',
        controller: 'QuerypageCtrl'
      })
      .when('/getPage/:id', {
        templateUrl: 'views/getpage.html',
        controller: 'GetpageCtrl'
      })
      .otherwise({
        redirectTo: '/queryPage'
      });

    $httpProvider.responseInterceptors.push(['$q', '$window', function($q, $window) {
      return function(promise) {
        return promise.then(function(response) {
          return response;
        }, function(response) {
          if (response.status !== 0) {
            $window.alert('予期しないエラー stauts:' + response.status);
          }
          return $q.reject(response);
        });
      };
    }]);
  });

  (function() {
    var importMockJs = function(jsfile) {
      if (location.hostname === 'mockhost') {
        /*jslint evil: true */
        document.write('<script type="text/javascript" src="' + jsfile + '"></script>');
      }
  };
  
  importMockJs('bower_components/angular-mocks/angular-mocks.js');
  importMockJs('scripts/app-mock.js');
  })();
