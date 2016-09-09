myApp.config(function ($stateProvider, $urlRouterProvider) {
	$urlRouterProvider.otherwise("/login");
	$stateProvider
	.state('login', {
		url : "/login",
		templateUrl : "/assets/templates/login-temp.html",
		controller : "LoginCtr"
	})
	.state('error', {
		url : "/error",
		templateUrl : "/assets/templates/error.html"
	})
	.state('admin', {
		url : "/admin",
		templateUrl : "/assets/templates/admin.html"
	})
	.state('admin.devices', {
		url : "/devices",
		templateUrl : "/assets/templates/devices.html",
		controller : "DevicesCtr"
	})
	.state('admin.site', {
		url : "/site",
		templateUrl : "/assets/templates/site.html",
		controller : "SiteCtr"
	})
	.state('admin.graph', {
		url : "/graph",
		templateUrl : "/assets/templates/graph.html",
		controller : "GraphCtr"
	});
});
