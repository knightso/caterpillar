'use strict';

/**
 * @ngdoc function
 * @name caterpillarApp.controller:PutpageCtrl
 * @description
 * # PutpageCtrl
 * Controller of the caterpillarApp
 */
angular.module('caterpillarApp')
  .controller('PutpageCtrl', function ($scope, Page, Property, $routeParams, $timeout, $http, method) {

    $scope.method = method;
    
    $scope.page = {};
    $scope.page.properties = {};
    
    if (method === 'PUT') {
      $scope.page = Page.get({
        id : $routeParams.id ? $routeParams.id : null
      });
    } else {
      $scope.page = {};
      $scope.page.properties = {};
    }

    // load leaves into pulldown-list.
    $scope.leaves = Page.getLeaves();
    $scope.leaves.$promise.then(function(leaves){
      // set default-value for update Page.
      if (method === 'PUT') {
        $scope.page = Page.get({
          id : $routeParams.id ? $routeParams.id : null
        });
        $scope.page.$promise.then(function(page) {
          $scope.originalPage = angular.copy(page);

          for (var i = 0; i < leaves.length; i++) {
            var leaf = leaves[i];
            if (leaf.name === page.leaf) {
              $scope.leaf = leaf;
              break;
            }
          }

          // load global property values.
          $scope.loadProperties();
        });
      }
    });

    var loadGlobalProperty = function(name) {
      Property.get({name : name}, function(property) {
        $scope.page.properties[name] = property.value;
      });
    };

    $scope.loadProperties = function() {
      if ($scope.leaf.wormholes !== undefined) {
        for (var i = 0; i < $scope.leaf.wormholes.length; i++) {
          var wormhole = $scope.leaf.wormholes[i];
          // TODO make it constants
          if (wormhole.type !== 'PROPERTY') {
            continue;
          }
          if (wormhole.global) {
            // TODO batch get
            loadGlobalProperty(wormhole.name);
          } else {
            if ($scope.originalPage !== undefined && $scope.originalPage.properties !== undefined) {
              $scope.page.properties[wormhole.name] = $scope.originalPage.properties[wormhole.name];
            }
          }
        }
      }
    };
  
    $scope.showProperties = function(wormhole) {
      return wormhole !== undefined && wormhole.type === 'PROPERTY';
    };

    $scope.validateNotAllNumbers = function() {
      if ($scope.page.alias !== undefined && $scope.page.alias !== '' && $scope.page.alias.match(/^[0-9]*$/)) {
        $scope.updatePageForm.alias.$error.allnumbers = true;
      } else {
        $scope.updatePageForm.alias.$error.allnumbers = false;
      }
    };

    $scope.submit = function() {

      $scope.alerts = [];

      if (!$scope.updatePageForm.$valid || $scope.updatePageForm.alias.$error.allnumbers) {
        $scope.alerts.push({type: 'danger', msg: '入力に不備がある為保存に失敗しました。\n各項目を見直して下さい。'});
        return;
      }
      
      $scope.page.leaf = $scope.leaf.name;
      for (var i = 0; i < $scope.leaf.wormholes.length; i++) {
        var wormhole = $scope.leaf.wormholes[i];
        if (wormhole.type === 'PROPERTY') {
          var property = $scope.page.properties[wormhole.name];
          if (property === undefined) {
            $scope.page.properties[wormhole.name] = '';
          }
        }
      }
      
      var start = +new Date();
      var tout = 100;

      var doSave = function(actualF) {
        actualF($scope.page,
          function() {
            var tat = +new Date() - start;
              $timeout(function() {
                $scope.alerts.push({type: 'success', msg: '保存に成功しました。'});
              }, (tat >= tout ? 0 : tout - tat));
            },
            function() {
              var tat = +new Date() - start;
              $timeout(function() {
                $scope.alerts.push({type: 'danger', msg: '保存に失敗しました。'});
              }, (tat >= tout ? 0 : tout - tat));
            }
        ); 
      };

      if ($scope.page.id === undefined) {
        doSave(Page.register);
      } else {
        doSave(Page.update);
      }
    };

    $scope.closeAlert = function(index) {
      $scope.alerts.splice(index, 1);
    };
  });
