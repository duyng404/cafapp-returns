var socket;
var socketToken = "";

window.onload = async function() {
	// send a get request to obtain socket token
	try {
		res = await getInfo();
		console.log(res);
		socketToken = res;
	} catch (error) {
		errorMessage("", "Some error happened while connecting:" + error + ". Maybe try reloading?");
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
		errorMessage("", "Cannot connect to TrackerBot T_T");
	}

	socket = io("/", {
		"path":"/socket/",
		"reconnectionAttempts":"5",
		"reconnectionDelay":"2000",
	})

	socket.on("connect", function () {
		console.log("connected to server from client");
		newMessage("", "TrackerBot online!", 'info');
	});
	socket.on("disconnect", function () {
		console.log("disconnected from server");
	});

	socket.on("newMessage", function (message) {
		console.log("message received from server")
		$("#chat-window").append(`<div class="card card-body bg-light"><span class="font-weight-bold">${message.user}</span>  ${message.content}</div>`);
	});
}

$("#main-input").on("submit", function (e) {
    e.preventDefault();
    var content = $("#main-input-text").val().trim();
    // socket.emit("createMessage", {
    //     user: "Lam",
    //     content: content
    // });

	console.log(content);
	$("#main-input-text").val("");
});

// newMessage
function newMessage(title, content, type) {
	var contentHTML = markdown.toHTML(content);
	if (!title) title = "";
	var specialBadge = "";
	if (type === "error") specialBadge = '<span class="badge badge-danger">Error</span>';
	if (type === "info") specialBadge = '<span class="badge badge-info">System</span>';
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
}

function errorMessage(title, error) {
	newMessage(title, error, 'error');
}
