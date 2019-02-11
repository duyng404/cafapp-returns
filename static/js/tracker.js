var socket;
var socketToken = "";

window.onload = async function() {
	// send a get request to obtain socket token
	try {
		res = await getInfo();
		console.log(res);
		socketToken = res;
	} catch (error) {
		errorMessage("", "Unable to connect to TrackerBot:" + error + ". Maybe try reloading?");
	}

	// init socket connection
	initSocket();
}

function getInfo() {
	return new Promise((resolve, reject) => {
		$.ajax({
			type: "GET",
			url: "/api/my-info",
			dataType: "json",
			success: function(data) {
				resolve(data.socket_token);
			},
			error: function (req, status, error) {
				console.log("get info failed with status",status,"and error",error);
				if (status === "error") reject(error)
				else reject(status);
			}
		})
	})
}

function initSocket() {
	if (socketToken === '') {
		console.log('token is blank');
		errorMessage("", "Unable to connect to TrackerBot: server error.");
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
				newMessage("", "TrackerBot online!", 'info');
			}
			else {
				console.log('invalid response from server');
				errorMessage("", "Unable to connect to TrackerBot: server error.");
				socket.close();
			}
		})
		// newMessage("", "TrackerBot online!", 'info');
	});

	socket.on('connect_error', function(error) {
		console.log('error connecting to socket:',error);
		errorMessage("", "Unable to connect to Trackerbot.");
	})
	socket.on('connect_timeout', function(timeout) {
		console.log('timed out when connecting to socket:',timeout);
		errorMessage("", "Unable to connect to Trackerbot: Timed Out.");
	})
	socket.on('reconnect_attempt', function(number) {
		console.log('attempting to connect #', number);
		errorMessage("", "Trying to reconnect...");
	})
	socket.on('disconnect', function (reason) {
		console.log('disconnected from server with reason', reason);
		errorMessage("", "Disconnected.");
	});

	socket.on("chatbot-response", function(response) {
		newMessage("", response, "chatbot");
	})
}

$("#main-input").on("submit", function (e) {
    e.preventDefault();
	var content = $("#main-input-text").val().trim();
	newMessage("You", content);
    socket.emit("chatbot-request", content, function(response){
		if (response === 'okbro') {
			console.log('okbro received');
			return
		} else {
			console.log('no okbro?');
			onConnectErrorFunc('error');
			setTimeout(() => {socket.open()},2000);
		}
	});
	$("#main-input-text").val("");
});

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

function updateScroll(){
	var element = document.getElementById("chat-scroll");
    element.scrollTop = element.scrollHeight;
}
