(function() {
	// jsonのキー名。
	var ID_KEY = 'id';
	var NAME_KEY = 'name';

	var app = angular.module('pageSelectorApp', [
		'pageSelectorControllers'
	]);

	var controllers = angular.module('pageSelectorControllers', []);
	controllers.controller('pageListCtrl', ['$scope', '$http', function($scope, $http) {
		var editor = window.top.tinymce.EditorManager.activeEditor;
		var params = editor.windowManager.getParams();

		$scope.pages = [];
		$http.get(params['requestUrl']).success(function(data) {
			console.log(data);
			$scope.pages = data;
		});	

		$scope.select = function(index) {
			var page = $scope.pages[index];
			var url = "caterpillar://" + page[ID_KEY]
			var name = page[NAME_KEY];

			var callback = params['callback'];
			callback(url, {text: name, title: name});
			editor.windowManager.close();
		}
	}]);
})();