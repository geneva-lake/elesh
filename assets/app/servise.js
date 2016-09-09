myApp.service('PagerService', ["$http", function ($http) {
			var service = {};
			service.GetPager = GetPager;
			return service;
			function GetPager(totalItems, currentPage, pageSize) {
				currentPage = currentPage || 1;
				pageSize = pageSize || 10;
				var totalPages = Math.ceil(totalItems / pageSize);
				var startPage,
				endPage;
				if (totalPages <= 10) {
					startPage = 1;
					endPage = totalPages;
				} else {
					if (currentPage <= 6) {
						startPage = 1;
						endPage = 10;
					} else if (currentPage + 4 >= totalPages) {
						startPage = totalPages - 9;
						endPage = totalPages;
					} else {
						startPage = currentPage - 5;
						endPage = currentPage + 4;
					}
				}
				var startIndex = (currentPage - 1) * pageSize;
				var endIndex = startIndex + pageSize;
				var pages = []; //range(startPage, endPage + 1);
				for (var i = 0; i < endPage; i++) {
					pages[i] = i + 1;
				}
				return {
					totalItems : totalItems,
					currentPage : currentPage,
					pageSize : pageSize,
					totalPages : totalPages,
					startPage : startPage,
					endPage : endPage,
					startIndex : startIndex,
					endIndex : endIndex,
					pages : pages
				};
			};
		}
	]);
