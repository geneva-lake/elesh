myApp.controller('GraphCtr', ['GraphModel', '$scope', '$location', function (GraphModel, $scope, $location) {
			var GraphCtr = this;

			GraphCtr.begin = new Date();
			GraphCtr.end = new Date();
			GraphCtr.begin.setDate(GraphCtr.begin.getDate() - 7);
			var riceData = {
				datasets :
				[{
						lineTension : 0,
						fillColor : "rgba(172,194,132,0.4)",
						strokeColor : "#ACC26D",
						pointColor : "#fff",
						pointStrokeColor : "#9DB86D",
					}
				]
			}
			$scope.dt = GraphCtr.begin;
			$scope.dt2 = GraphCtr.end;
			$scope.inlineOptions = {
				customClass : getDayClass,
				minDate : new Date(),
				showWeeks : true
			};
			$scope.dateOptions = {
				dateDisabled : disabled,
				formatYear : 'yy',
				maxDate : new Date(2020, 5, 22),
				minDate : new Date(),
				startingDay : 1
			};
			function disabled(data) {
				var date = data.date,
				mode = data.mode;
				return mode === 'day' && (date > new Date());
			};
			$scope.toggleMin = function () {
				$scope.inlineOptions.minDate = $scope.inlineOptions.minDate ? null : new Date();
				$scope.dateOptions.minDate = $scope.inlineOptions.minDate;
			};
			$scope.toggleMin();
			$scope.open1 = function () {
				$scope.popup1.opened = true;
			};
			$scope.open2 = function () {
				$scope.popup2.opened = true;
			};
			$scope.formats = ['dd-MMMM-yyyy', 'yyyy/MM/dd', 'dd.MM.yyyy', 'shortDate'];
			$scope.format = $scope.formats[0];
			$scope.altInputFormats = ['M!/d!/yyyy'];

			$scope.popup1 = {
				opened : false
			};

			$scope.popup2 = {
				opened : false
			};

			function getDayClass(data) {
				var date = data.date,
				mode = data.mode;
				if (mode === 'day') {
					var dayToCheck = new Date(date).setHours(0, 0, 0, 0);

					for (var i = 0; i < $scope.events.length; i++) {
						var currentDay = new Date($scope.events[i].date).setHours(0, 0, 0, 0);

						if (dayToCheck === currentDay) {
							return $scope.events[i].status;
						}
					}
				}
				return '';
			}
			GraphCtr.dataChng = function (n, o) {
				GraphModel.all($scope.dt, $scope.dt2).success(function (data) {

					riceData.labels = data.date;
					riceData.datasets[0].data = data.count;

					if ($scope.chart) {
						$scope.chart.destroy();
					}
					var rice = document.getElementById('my-chart').getContext('2d');
					$scope.chart = new Chart(rice).Line(riceData);

				}).error(function (data, header, config) {
					if (header === 401) {
						$location.path('/login');
					} else {
						$location.path('/error');
					}
				});
			};

			$scope.$watch('dt', GraphCtr.dataChng);
			$scope.$watch('dt2', GraphCtr.dataChng);
		}
	]);
