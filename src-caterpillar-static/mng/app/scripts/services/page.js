'use strict';

/**
 * @ngdoc service
 * @name caterpillarApp.page
 * @description
 * # page
 * Service in the caterpillarApp.
 */
angular.module('caterpillarApp')
  .service('Page', function Page($resource) {
    return $resource('/caterpillar/api/pages/:id', {id:'@id'}, {
      query: {method:'GET', url:'/caterpillar/api/pages', isArray:true},
      update: {method:'PUT'},
      register: {method:'POST'},
      getLeaves: {method:'GET', url:'/caterpillar/api/leaves', isArray:true}
    });
  });
