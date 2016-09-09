myApp.controller('SiteCtr', ['$scope', '$http', '$location', function ($scope, $http, $location) {
			$http.get("http://localhost:3000/admin/site")
			.success(
				function (data) {
				$scope.siteUrl = data.url;
				$scope.siteDate = data.date;
			}).error(function (data, header, config) {
				if (header === 401) {
					$location.path('/login');
				} else {
					$location.path('/error');
				}
			});
			$scope.updateSiteUrl = function (data) {
				$scope.siteDate = new Date();
				$http.post("http://localhost:3000/admin/site", {
					url : data,
					date : $scope.siteDate
				});
			};

		}
	]);
