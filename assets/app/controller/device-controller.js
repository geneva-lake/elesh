myApp.controller('DevicesCtr', ['DevicesModel', 'PagerService', '$scope', '$location', 
	function (DevicesModel, PagerService, $scope, $location) {
			var devicesCtr = this;
			$scope.orderId = 1;
			$scope.orderDate = 1;
			devicesCtr.length = 5;
			$scope.setPage = function (page) {
				if (page < 1 || page > devicesCtr.pager.totalPages) {
					return;
				}
				devicesCtr.pager = PagerService.GetPager($scope.total, page, devicesCtr.length);
				$scope.currentPage = page;
				DevicesModel.all(page, devicesCtr.length, "device-id", -1).success(function (data) {
					$scope.ScopeDevices = data.devices;
					$scope.length = 5;
					$scope.total = data.Total;
					var firststep = 1;
					$scope.totalPages = devicesCtr.pager.totalPages;
					$scope.pages = devicesCtr.pager.pages;

				}).error(function (data, header, config) {
					if (header === 401) {
						$location.path('/login');
					} else {
						$location.path('/error');
					}
				});
			};
			function initController() {
				devicesCtr.pager = {};
				var firststep = 1;
				$scope.pages = devicesCtr.pager.pages;
				$scope.totalPages = devicesCtr.pager.totalPages;
				$scope.setPage(1);
			}
			initController();
			$scope.sort_by = function (filter) {
				DevicesModel.all($scope.currentPage, 5, filter, $scope.orderId)
				.success(function (data) {
					$scope.ScopeDevices = data.devices;
					$scope.pages = devicesCtr.pager.pages;
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
