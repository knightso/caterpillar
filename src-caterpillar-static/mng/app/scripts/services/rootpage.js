'use strict';

/**
 * @ngdoc service
 * @name caterpillarApp.Rootpage
 * @description
 * # Rootpage
 * Service in the caterpillarApp.
 */
angular.module('caterpillarApp')
  .service('Rootpage', function Rootpage($resource) {
    return $resource('/caterpillar/api/rootPage', {}, {
      update: {method:'PUT'}
    });
  });
