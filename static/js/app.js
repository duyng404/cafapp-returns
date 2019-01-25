//CLIENT

const socket = io.connect('http://localhost:7000',{
    "path": "/socket/" 
});
socket.on('connect', function () {
    console.log('connected to server from client');
});
socket.on('disconnect', function () {
    console.log('disconnected from server');
});

socket.on('newMessage', function (message) {
    console.log("message received from server")
    $("#chat-window").append(`<div class="card card-body bg-light"><span class="font-weight-bold">${message.user}</span>  ${message.content}</div>`);
});

$("#send").on("click", function (e) {
    e.preventDefault();
    var content = $("#msg-content").val().trim();
    socket.emit("createMessage", {
        user: "Lam",
        content: content
    });
    $("#msg-content").val('');
});