'use strict';

/**
 * @ngdoc service
 * @name caterpillarApp.property
 * @description
 * # property
 * Service in the caterpillarApp.
 */
angular.module('caterpillarApp')
  .service('Property', function Property($resource) {
    return $resource('/caterpillar/api/property/:name', {id:'@name'}, {
      query: {method:'GET', url:'/caterpillar/api/properties', isArray:true},
      update: {method:'PUT'},
      register: {method:'POST'}
    });
  });
