myApp.service('DevicesModel', ["$http", function ($http) {
			var service = this;
			service.all = function (skip, limit, filter, order) {
				return $http.get("http://localhost:3000/admin/devices/" + skip + "/" + limit + "/" + filter + "/" + order)
				.success(function (data) {
					return data;
				}).
				error(function (data, header, config) {
					return null;
				});

			};
		}
	]);

myApp.service('GraphModel', ["$http", function ($http) {
			var service = this;
			service.all = function (begin, end) {
				var postDate = {};
				postDate.begin = begin;
				postDate.end = end;
				return $http.post("http://localhost:3000/admin/use-count", postDate)
				.success(function (data) {
					return data;
				}).
				error(function (data, headers, config) {
					return null;
				});
			};
		}
	]);
