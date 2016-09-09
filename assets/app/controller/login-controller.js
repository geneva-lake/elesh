myApp.controller('LoginCtr', ['$scope', '$location', '$http', function ($scope, $location, $http) {
			var loginCtr = this;
			$scope.login = function () {
				postData = {
					login : $scope.formLogin,
					password : $scope.formPassword
				};
				$http.post('http://localhost:3000/web', postData).
				success(function (data) {
					$location.path('/admin');
				}).error(function (data, header, config) {
					if (header === 401) {
						$location.path('/login');
					} else {
						$location.path('/error');
					}
				});
			};

		}
	]);
