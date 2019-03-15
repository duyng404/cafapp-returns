var socket;
var lastPolled = 0;
var socketToken = "";
var hasActiveOrders = false;

window.onload = async function() {
	setListeners();

	// send a get request to obtain socket token
	try {
		res = await getInfo();
		console.log(res);
		socketToken = res.socket_token;
		hasActiveOrders = res.has_active_orders;
	} catch (error) {
		showAlert("Cannot connect to server:"+error+". Please try again.")
	}

	if (hasActiveOrders) {
		// init UI
		initTrackerUI();
		// set up long polling
		initLongPolling();
		// init socket connection
		initSocket();
	} else {
		$("#orderTrackerContainer").html("You don't have any active orders to track");
	}
}

function getInfo() {
	return new Promise((resolve, reject) => {
		$.ajax({
			type: "GET",
			url: "/api/my-info",
			dataType: "json",
			success: function(data) {
				console.log("getInfo() success!");
				resolve(data);
			},
			error: function (req, status, error) {
				console.log("get info failed with status",status,"and error",error);
				if (status === "error") reject(error)
				else reject(status);
			}
		})
	})
}

function initTrackerUI() {
	$("#orderTrackerContainer").html(`
		<div class="d-flex justify-content-between">
			<small id="trackerConnectivity"></small>
			<small class="text-muted">Last Updated: <span id="trackerLastUpdated"></span></small>
		</div>
		<div class="progress">
			<div id="trackerProgressBar" class="progress-bar pulsing-background" role="progressbar" style="width: 0%" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100"></div>
		</div>
		<div id="trackerContent" class="mt-3">
			<p>loading...</p>
		</div>
	`);
	setLastUpdated(moment());
}

function initLongPolling() {
	this.interval = setInterval(fetchActiveOrders, 1000);
}

function initSocket() {
	if (socketToken === '') {
		console.log('token is blank');
		setConnectivity("Connection Error","text-danger");
	}

	socket = io("/", {
		"path":"/socket/",
		"reconnectionAttempts":"5",
		"reconnectionDelay":"2000",
	})

	socket.on("connect", function () {
		socket.emit('register-user', socketToken, function(response){
			if (response === "okbro") {
				console.log('socket connected!');
				setConnectivity("Connected","text-muted");
			}
			else {
				console.log('invalid response from server');
				setConnectivity("Connection Error","text-danger");
				socket.close();
			}
		})
	});

	socket.on('connect_error', function(error) {
		console.log('error connecting to socket:',error);
		setConnectivity("Connection Error","text-danger");
	})
	socket.on('connect_timeout', function(timeout) {
		console.log('timed out when connecting to socket:',timeout);
		setConnectivity("Connection Timed Out","text-danger");
	})
	socket.on('reconnect_attempt', function(number) {
		console.log('attempting to connect #', number);
		setConnectivity("Reconnecting","text-warning");
	})
	socket.on('disconnect', function (reason) {
		console.log('disconnected from server with reason', reason);
		setConnectivity("Disconnected","text-danger");
	});

	socket.on("active-orders-update", function(data) {
		console.log(data);
		fetchActiveOrders(true);
	})
}

function fetchActiveOrders(force = false) {
	var now = moment().unix();
	if (force || now - lastPolled > 20) {
		lastPolled = now;
		$.ajax({
			type: "GET",
			url: "/api/view-active-orders",
			contentType: "application/json",
			dataType: "json",
			success: function(data) {
				setActiveOrders(data);
			},
			error: function (req, status, error) {
				console.log("view active orders failed with status",status,"and error",error);
				showAlert("Error happened while getting active orders. Maybe try refreshing?")
			}
		})
	}
}

function getFriendlyIDFromTag(s) {
	var ss = s.split('-');
	return ss[ss.length-1];
}

function setActiveOrders(orders) {
	var progress = 0;
	var totalProgress = 0;
	$("#trackerContent").empty();
	if (orders.length && orders.length > 0) {
		for (let o of orders) {
			var status
			switch (o.status_code) {
				case 20:
					status="Placed";
					progress += 20;
					break;
				case 40:
					status="Prepping";
					progress += 40;
					break;
				case 50:
					status="Out for delivery";
					progress += 60;
					break;
				case 51:
					status="Arrived at door";
					progress += 80;
					showAlert(`An order is coming in just a few seconds. Please come and pick it up at ${o.destination.pickup_location}!`, "info")
					break;
				case 60:
					status="Delivered";
					progress += 100;
					break;
			}
			totalProgress += 100
			$("#trackerContent").append(`
				<p>
					Your order #${getFriendlyIDFromTag(o.tag)} going to ${o.destination.name} is ${status} (updated ${moment(o.updated_at).format("h:mm a")})
					<br />Pick Up Location: ${o.destination.pickup_location}
				</p>
			`);
		}
	}
	var percent = Math.round(progress/totalProgress*100);
	$("#trackerProgressBar").css("width",percent + "%");
	$("#trackerProgressBar").prop("aria-valuenow", percent + "");
	setLastUpdated(moment());
}

// $("#main-input").on("submit", function (e) {
//     e.preventDefault();
// 	var content = $("#main-input-text").val().trim();
// 	newMessage("You", content);
//     socket.emit("chatbot-request", content, function(response){
// 		if (response === 'okbro') {
// 			console.log('okbro received');
// 			return
// 		} else {
// 			console.log('no okbro?');
// 			onConnectErrorFunc('error');
// 			setTimeout(() => {socket.open()},2000);
// 		}
// 	});
// 	$("#main-input-text").val("");
// });

// newMessage
function newMessage(title, content, type) {
	var contentHTML = markdown.toHTML(content);
	if (!title) title = "";
	var specialBadge = "";
	if (type === "error") specialBadge = '<span class="badge badge-danger">Error</span>';
	if (type === "info") specialBadge = '<span class="badge badge-info">System</span>';
	if (type === "chatbot") specialBadge = '<span class="badge badge-success">TrackerBot</span>';
	var timestamp = moment().format("h:mm a");
	$('\
	<div class="bg-white rounded shadow-sm my-2">\
		<div class="media pt-3">\
			<div class="media-body px-3 mb-0">\
				'+specialBadge+'\
				<strong>'+title+'</strong>&nbsp;<small class="font-weight-light">\
				'+timestamp+'\
				</small>\
				'+contentHTML+'\
			</div>\
		</div>\
	</div>').appendTo("#chat-window");
	updateScroll();
}

function errorMessage(title, error) {
	newMessage(title, error, 'error');
}

// function updateScroll(){
// 	var element = document.getElementById("chat-scroll");
//     element.scrollTop = element.scrollHeight;
// }

function setConnectivity(status, color) {
	$("#trackerConnectivity").html(`
		<span class="${color}">${status}</span>
	`);
}

function setLastUpdated(time) {
	$("#trackerLastUpdated").html(moment(time).format("h:mm a"));
}

// NON-TRACKER STUFF

function showAlert(msg, color) {
	$("#toast-area").html("");
	if (!msg) return
	if (!color) color = "primary"
	$("#toast-area").prepend(`
		<div class="alert alert-${color} alert-dismissible fade show" role="alert">
			${msg}
			<button type="button" class="close" data-dismiss="alert" aria-label="Close">
				<span aria-hidden="true">&times;</span>
			</button>
		</div>
	`);
}

function clearPhone() {
	$("#phone").html("");
}

function showPhoneEdit() {
	phoneNum = $("#current-phone-num").val();
	$("#phone").prepend(`
		<div class="input-group input-group-sm">
			<input class="form-control font-codecard" type="text" id="phone-number" data-inputmask="'mask': '(999)-999-9999'"/>
			<div class="input-group-append">
				<button class="btn btn-outline-ca-yellow4" type="button" id="phone-number-submit">Update</button>
			</div>
		</div>
		<a class="px-3" href="#" role="button" id="phone-number-cancel">cancel</a>
	`);
	$(":input").inputmask();
	$("#phone-number").focus().val(phoneNum);
	setListeners();
}

function showPhoneNormal() {
	phoneNum = $("#current-phone-num").val();
	$("#phone").prepend(`
		<span id="phone-number">${phoneNum}</span>
		<a class="px-3" href="#" role="button" id="phone-number-edit">edit</a>
	`);
	setListeners();
}

function setListeners() {
	$("#phone-number-edit").on("click", function () {
		clearPhone();
		showPhoneEdit();
		return false;
	});

	$("#phone-number-submit").on("click", function() {
		phoneNum = $("#phone-number").val();
		$.ajax({
			method: 'POST',
			url: '/api/edit-phone',
			contentType: 'application/json',
			dataType: 'json',
			data: JSON.stringify({phone: phoneNum}),
			success: function (res) {
				console.log(res);
				$("#current-phone-num").val(res);
				clearPhone();
				showPhoneNormal();
				showAlert("Phone change success!", "success")
			},
			error: function (err) {
				console.log(err);
				clearPhone();
				showPhoneNormal();
				showAlert("Oh no! Phone change failed! Please try again later.")
			}
		})
		return false;
	});

	$("#phone-number-cancel").on("click", function() {
		clearPhone();
		showPhoneNormal();
		return false;
	})
}
