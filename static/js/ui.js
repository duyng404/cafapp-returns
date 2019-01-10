function toggleNav(){
	$("#nav-mainmenu").toggleClass("show");
}

function toggleUsr(){
	$("#nav-usermenu").toggleClass("show");
}

function scrollToInfo(){
	document.querySelector('#howw').scrollIntoView({ 
		block: "start",
		inline: "nearest",
		behavior: 'smooth' 
	});
}