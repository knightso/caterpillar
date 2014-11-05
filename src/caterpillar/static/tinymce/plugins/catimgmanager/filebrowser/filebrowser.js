(function() {
	var app = angular.module('imageManagerApp', [
		'imageManagerControllers'
	]);

	var controllers = angular.module('imageManagerControllers', []);
	controllers.controller('imageListCtrl', ['$scope', '$http', function($scope, $http) {
		var editor = window.top.tinymce.EditorManager.activeEditor;
		var params = editor.windowManager.getParams();

		$scope.images = [];
		$http.get(params['requestUrl']).success(function(data) {
			$scope.images = data['files'];
		});

		$scope.select = function(index) {
			var img = $scope.images[index];
			var callback = params['callback'];
			callback(img.url, {alt: img.filename});
			editor.windowManager.close();
		};

		var uploadWindow;
		$scope.upload = function() {
			uploadWindow = editor.windowManager.open({
				title: '画像のアップロード',
				url: '/caterpillar/static/tinymce/plugins/catimgmanager/dialog.html',
				width: 300,
				height: 240,
				buttons: [
					{text: 'キャンセル', onclick: 'close'}
				]
			});
		}

		window.parent.uploadSuccess = function() {
			$http.get(params['requestUrl']).success(function(data) {
					$scope.images = data['files'];
					console.log($scope.images.length);
					uploadWindow.close();
				});
		}
	}]);
})();